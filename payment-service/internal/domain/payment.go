package domain

import (
	"context"
)

type Payment struct {
	ID            string `json:"id"`
	OrderID       string `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	Amount        int64  `json:"amount"` // Amount in cents
	Status        string `json:"status"` // "Authorized", "Declined"
}

// PaymentRepository interface for interacting with underlying datastore.
type PaymentRepository interface {
	Store(ctx context.Context, p *Payment) error
	GetByOrderID(ctx context.Context, orderID string) (*Payment, error)
}

// PaymentUseCase interface represents business logic.
type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, orderID string, amount int64) (*Payment, error)
	GetPaymentStatus(ctx context.Context, orderID string) (*Payment, error)
}
