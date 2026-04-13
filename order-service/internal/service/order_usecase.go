package service

import (
	"context"
	"errors"
	"order-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ErrServiceUnavailable struct {
	Message string
}

func (e *ErrServiceUnavailable) Error() string {
	return e.Message
}

type orderUseCase struct {
	repo          domain.OrderRepository
	paymentClient domain.PaymentClient
	broker        *OrderBroker
}

func NewOrderUseCase(repo domain.OrderRepository, paymentClient domain.PaymentClient, broker *OrderBroker) domain.OrderUseCase {
	return &orderUseCase{
		repo:          repo,
		paymentClient: paymentClient,
		broker:        broker,
	}
}

func (u *orderUseCase) CreateOrder(ctx context.Context, customerID, itemName string, amount int64, idempotencyKey string) (*domain.Order, error) {
	if idempotencyKey != "" {
		existingOrder, err := u.repo.GetByIdempotencyKey(ctx, idempotencyKey)
		if err == nil && existingOrder != nil {
			return existingOrder, nil
		}
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	order := &domain.Order{
		ID:             uuid.New().String(),
		CustomerID:     customerID,
		ItemName:       itemName,
		Amount:         amount,
		Status:         "Pending",
		CreatedAt:      time.Now(),
		IdempotencyKey: idempotencyKey,
	}

	if err := u.repo.Store(ctx, order); err != nil {
		return nil, err
	}

	go func(orderID string, amount int64) {

		bgCtx := context.Background()

		time.Sleep(15 * time.Second) // pending status simulating

		

		_, err := u.paymentClient.AuthorizePayment(bgCtx, orderID, amount)
		if err != nil {
			_ = u.repo.UpdateStatus(bgCtx, orderID, "Payment Failed")
			u.broker.Publish(orderID, "Payment Failed")
			return
		}

		_ = u.repo.UpdateStatus(bgCtx, orderID, "Paid")
		u.broker.Publish(orderID, "Paid")
	}(order.ID, order.Amount)

	// to give while Penging
	return order, nil

}

func (u *orderUseCase) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	o, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if o == nil {
		return nil, errors.New("order not found")
	}
	return o, nil
}

func (u *orderUseCase) CancelOrder(ctx context.Context, id string) error {
	o, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if o == nil {
		return errors.New("order not found")
	}

	if o.Status == "Paid" {
		return errors.New("cannot cancel a Paid order")
	}
	if o.Status != "Pending" {
		return errors.New("only Pending orders can be cancelled")
	}

	err = u.repo.UpdateStatus(ctx, id, "Cancelled")
	if err == nil {
		u.broker.Publish(id, "Cancelled")
	}

	return err
}

func isTimeoutOrUnavailable(err error) bool {
	return err != nil && (err.Error() == "payment service unavailable" || len(err.Error()) > 27 && err.Error()[:27] == "payment service unavailable")
}

func (u *orderUseCase) GetRevenueByCustomerID(ctx context.Context, customerID string) (*domain.CustomerRevenue, error) {
	if customerID == "" {
		return nil, errors.New("customer_id is required")
	}
	return u.repo.GetRevenueByCustomerID(ctx, customerID)
}
