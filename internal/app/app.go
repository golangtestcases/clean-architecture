package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/golangtestcases/subscribe-service/docs"
	"github.com/golangtestcases/subscribe-service/internal/app/handlers/create_subscription_handler"
	"github.com/golangtestcases/subscribe-service/internal/app/handlers/delete_subscription_handler"
	"github.com/golangtestcases/subscribe-service/internal/app/handlers/get_cost_handler"
	"github.com/golangtestcases/subscribe-service/internal/app/handlers/get_subscription_handler"
	"github.com/golangtestcases/subscribe-service/internal/app/handlers/list_subscriptions_handler"
	"github.com/golangtestcases/subscribe-service/internal/app/handlers/update_subscription_handler"
	"github.com/golangtestcases/subscribe-service/internal/domain/subscription/repository"
	"github.com/golangtestcases/subscribe-service/internal/domain/subscription/service"
	"github.com/golangtestcases/subscribe-service/internal/infra/config"
	"github.com/golangtestcases/subscribe-service/internal/infra/http/middlewares"
)

type App struct {
	config *config.Config
	server http.Server
	db     *sql.DB
	logger *slog.Logger
}

func NewApp(configPath string) (*App, error) {
	configImpl, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	logger := setupLogger(configImpl.Logger.Level)

	db, err := setupDatabase(configImpl.Database)
	if err != nil {
		return nil, fmt.Errorf("setupDatabase: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("runMigrations: %w", err)
	}

	app := &App{
		config: configImpl,
		db:     db,
		logger: logger,
	}

	app.server.Handler = bootstrapHandler(db, logger)

	return app, nil
}

func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)
	app.logger.Info("starting server", "address", address)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return app.server.Serve(l)
}

func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	return slog.New(handler)
}

func setupDatabase(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func bootstrapHandler(db *sql.DB, logger *slog.Logger) http.Handler {
	subscriptionRepository := repository.NewPostgreSQLRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepository)

	mx := http.NewServeMux()

	mx.Handle("POST /api/subscriptions", create_subscription_handler.NewCreateSubscriptionHandler(subscriptionService, logger))
	mx.Handle("GET /api/subscriptions/{id}", get_subscription_handler.NewGetSubscriptionHandler(subscriptionService, logger))
	mx.Handle("PUT /api/subscriptions/{id}", update_subscription_handler.NewUpdateSubscriptionHandler(subscriptionService, logger))
	mx.Handle("DELETE /api/subscriptions/{id}", delete_subscription_handler.NewDeleteSubscriptionHandler(subscriptionService, logger))
	mx.Handle("GET /api/subscriptions", list_subscriptions_handler.NewListSubscriptionsHandler(subscriptionService, logger))
	mx.Handle("GET /api/subscriptions/cost", get_cost_handler.NewGetCostHandler(subscriptionService, logger))

	// Swagger
	mx.Handle("GET /swagger/", httpSwagger.WrapHandler)

	middleware := middlewares.NewTimerMiddleware(mx)

	return middleware
}
