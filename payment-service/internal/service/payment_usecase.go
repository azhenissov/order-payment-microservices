package service

import (
	"context"
	"encoding/json" // Добавили для маршалинга
	"errors"
	"log" // Добавили для логирования ошибок отправки
	"payment-service/internal/domain"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type paymentUseCase struct {
	repo     domain.PaymentRepository
	rabbitCh *amqp.Channel 
}

type OrderEvent struct {
	OrderID       string  `json:"order_id"`
	Amount        float64 `json:"amount"`
	CustomerEmail string  `json:"customer_email"`
	Status        string  `json:"status"`
}

//  Обновили конструктор, чтобы принимать rabbitCh
func NewPaymentUseCase(repo domain.PaymentRepository, rabbitCh *amqp.Channel) domain.PaymentUseCase {
	return &paymentUseCase{
		repo:     repo,
		rabbitCh: rabbitCh,
	}
}

func (u *paymentUseCase) ProcessPayment(ctx context.Context, orderID string, amount int64) (*domain.Payment, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	status := "Authorized"
	if amount > 100000 {
		status = "Declined"
	}

	payment := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
		Status:        status,
	}

	// Сначала сохраняем в БД
	err := u.repo.Store(ctx, payment)
	if err != nil {
		return nil, err
	}

	//  Отправляем событие в RabbitMQ после успешной записи в БД
	event := OrderEvent{
		OrderID:       payment.OrderID,
		Amount:        float64(payment.Amount),
		CustomerEmail: "customer@example.com", 
		Status:        payment.Status,
	}

	body, err := json.Marshal(event)
	if err == nil {
		err = u.rabbitCh.PublishWithContext(ctx,
			"",                  // exchange
			"payment.completed", // routing key
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent, // Гарантия надежности
				Body:         body,
			})
		if err != nil {
			// Логируем ошибку, но не прерываем выполнение (платеж-то прошел)
			log.Printf("Failed to publish payment event: %v", err)
		} else {
			log.Printf("Event published for OrderID: %s", payment.OrderID)
		}
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

func (uc *paymentUseCase) ListPayments(ctx context.Context, min, max int64) ([]*domain.Payment, error) {
	if min > 0 && max > 0 && min > max {
		return nil, errors.New("Min cannot be less than Max ")
	}
	return uc.repo.FindByAmountRange(ctx, min, max)
}