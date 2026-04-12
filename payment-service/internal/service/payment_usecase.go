package service

import (
	"context"
	"errors"
	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type paymentUseCase struct {
	repo domain.PaymentRepository
}

func NewPaymentUseCase(repo domain.PaymentRepository) domain.PaymentUseCase {
	return &paymentUseCase{repo: repo}
}

func (u *paymentUseCase) ProcessPayment(ctx context.Context, orderID string, amount int64) (*domain.Payment, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	status := "Authorized"
	if amount > 100000 { // 1000 units (e.g. $1000) -> Declined
		status = "Declined"
	}

	payment := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
		Status:        status,
	}

	err := u.repo.Store(ctx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (u *paymentUseCase) GetPaymentStatus(ctx context.Context, orderID string) (*domain.Payment, error) {
	payment, err := u.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errors.New("payment not found for the given order id")
	}
	return payment, nil
}
