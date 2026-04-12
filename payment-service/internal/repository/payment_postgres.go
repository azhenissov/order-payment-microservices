package repository

import (
	"context"
	"database/sql"
	"payment-service/internal/domain"
)

type postgresPaymentRepository struct {
	db *sql.DB
}

// NewPostgresPaymentRepository creates a new postgres repository
func NewPostgresPaymentRepository(db *sql.DB) domain.PaymentRepository {
	return &postgresPaymentRepository{db: db}
}

func (r *postgresPaymentRepository) Store(ctx context.Context, p *domain.Payment) error {
	query := `INSERT INTO payments (id, order_id, transaction_id, amount, status) 
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query, p.ID, p.OrderID, p.TransactionID, p.Amount, p.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresPaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	query := `SELECT id, order_id, transaction_id, amount, status FROM payments WHERE order_id = $1`

	p := &domain.Payment{}
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}
