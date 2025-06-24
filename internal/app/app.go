package app

import (
	"fmt"
	"net"
	"net/http"

	"github.com/golangtestcases/clean-architecture/internal/app/handlers/create_entity_handler"
	"github.com/golangtestcases/clean-architecture/internal/app/handlers/get_entity_handler"
	"github.com/golangtestcases/clean-architecture/internal/domain/entities/repository"
	"github.com/golangtestcases/clean-architecture/internal/domain/entities/service"
	"github.com/golangtestcases/clean-architecture/internal/infra/config"
	"github.com/golangtestcases/clean-architecture/internal/infra/http/middlewares"
)

type App struct {
	config *config.Config
	server http.Server
}

func NewApp(configPath string) (*App, error) {
	configImpl, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{
		config: configImpl,
	}

	app.server.Handler = bootstrapHandler()

	return app, nil
}

func (app *App) ListenAndServe() error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return app.server.Serve(l)
}

func bootstrapHandler() http.Handler {
	entityRepository := repository.NewInMemoryRepository(100)
	entityService := service.NewEntityService(entityRepository)

	mx := http.NewServeMux()
	mx.Handle("POST /api/entities", create_entity_handler.NewCreateEntityHandler(entityService))
	mx.Handle("GET /api/entities/{id}", get_entity_handler.NewGetEntityHandler(entityService))

	middleware := middlewares.NewTimerMiddleware(mx)

	return middleware
}
