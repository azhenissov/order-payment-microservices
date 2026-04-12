package domain

import (
	"context"
	"time"
)

type Order struct {
	ID             string    `json:"id"`
	CustomerID     string    `json:"customer_id"`
	ItemName       string    `json:"item_name"`
	Amount         int64     `json:"amount"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	IdempotencyKey string    `json:"idempotency_key,omitempty"`
}

type CustomerRevenue struct {
	CustomerID  string `json:"customer_id"`
	TotalAmount int64  `json:"total_amount"`
	OrdersCount int64  `json:"orders_count"`
}

type OrderRepository interface {
	Store(ctx context.Context, o *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	GetRevenueByCustomerID(ctx context.Context, customerID string) (*CustomerRevenue, error)
}

type PaymentClient interface {
	AuthorizePayment(ctx context.Context, orderID string, amount int64) (string, error)
}

type OrderUseCase interface {
	CreateOrder(ctx context.Context, customerID, itemName string, amount int64, idempotencyKey string) (*Order, error)
	GetOrder(ctx context.Context, id string) (*Order, error)
	CancelOrder(ctx context.Context, id string) error
	GetRevenueByCustomerID(ctx context.Context, customerID string) (*CustomerRevenue, error)
}
