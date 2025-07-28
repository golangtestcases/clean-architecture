package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golangtestcases/subscribe-service/internal/domain/model"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewPostgreSQLRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

func (r *PostgreSQLRepository) CreateSubscription(ctx context.Context, subscription model.Subscription) (model.Subscription, error) {
	subscription.ID = uuid.New()
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()

	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		subscription.ID, subscription.ServiceName, subscription.Price,
		subscription.UserID, subscription.StartDate, subscription.EndDate,
		subscription.CreatedAt, subscription.UpdatedAt,
	)

	return subscription, err
}

func (r *PostgreSQLRepository) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (model.Subscription, error) {
	var subscription model.Subscription
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID, &subscription.ServiceName, &subscription.Price,
		&subscription.UserID, &subscription.StartDate, &subscription.EndDate,
		&subscription.CreatedAt, &subscription.UpdatedAt,
	)

	return subscription, err
}

func (r *PostgreSQLRepository) UpdateSubscription(ctx context.Context, subscription model.Subscription) error {
	subscription.UpdatedAt = time.Now()

	query := `
		UPDATE subscriptions 
		SET service_name = $2, price = $3, user_id = $4, start_date = $5, end_date = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		subscription.ID, subscription.ServiceName, subscription.Price,
		subscription.UserID, subscription.StartDate, subscription.EndDate,
		subscription.UpdatedAt,
	)

	return err
}

func (r *PostgreSQLRepository) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgreSQLRepository) ListSubscriptions(ctx context.Context, limit, offset int) ([]model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription
	for rows.Next() {
		var subscription model.Subscription
		err := rows.Scan(
			&subscription.ID, &subscription.ServiceName, &subscription.Price,
			&subscription.UserID, &subscription.StartDate, &subscription.EndDate,
			&subscription.CreatedAt, &subscription.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, rows.Err()
}

func (r *PostgreSQLRepository) GetTotalCost(ctx context.Context, filter model.CostFilter) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name ILIKE $%d", argIndex)
		args = append(args, "%"+*filter.ServiceName+"%")
		argIndex++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND start_date >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND (end_date IS NULL OR end_date <= $%d)", argIndex)
		args = append(args, *filter.EndDate)
	}

	var totalCost int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&totalCost)
	return totalCost, err
}
