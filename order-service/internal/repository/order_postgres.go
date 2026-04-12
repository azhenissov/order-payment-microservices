package repository

import (
	"context"
	"database/sql"
	"order-service/internal/domain"
)

type postgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) domain.OrderRepository {
	return &postgresOrderRepository{db: db}
}

func (r *postgresOrderRepository) Store(ctx context.Context, o *domain.Order) error {
	query := `INSERT INTO orders (id, customer_id, item_name, amount, status, created_at, idempotency_key) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, o.ID, o.CustomerID, o.ItemName, o.Amount, o.Status, o.CreatedAt, o.IdempotencyKey)
	return err
}

func (r *postgresOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `SELECT id, customer_id, item_name, amount, status, created_at, idempotency_key FROM orders WHERE id = $1`
	o := &domain.Order{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt, &o.IdempotencyKey,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	return o, err
}

func (r *postgresOrderRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Order, error) {
	if key == "" {
		return nil, nil
	}
	query := `SELECT id, customer_id, item_name, amount, status, created_at, idempotency_key FROM orders WHERE idempotency_key = $1`
	o := &domain.Order{}
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt, &o.IdempotencyKey,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	return o, err
}

func (r *postgresOrderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *postgresOrderRepository) GetRevenueByCustomerID(ctx context.Context, customerID string) (*domain.CustomerRevenue, error) {
	query := `SELECT COALESCE(SUM(amount), 0), COUNT(id) FROM orders WHERE customer_id = $1 and status = 'Paid'`
	var totalAmount int64
	var count int64
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&totalAmount, &count)
	if err != nil {
		return nil, err
	}
	return &domain.CustomerRevenue{
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		OrdersCount: count,
	}, nil
}//
