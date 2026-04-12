package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/internal/domain"
	"time"
)

type httpPaymentClient struct {
	client     *http.Client
	paymentURL string
}

func NewHTTPPaymentClient(paymentURL string) domain.PaymentClient {
	// Custom HTTP client with max 2s timeout
	return &httpPaymentClient{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
		paymentURL: paymentURL,
	}
}

type paymentRequest struct {
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount"`
}

type paymentResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
}

func (c *httpPaymentClient) AuthorizePayment(ctx context.Context, orderID string, amount int64) (string, error) {
	reqBody := paymentRequest{
		OrderID: orderID,
		Amount:  amount,
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.paymentURL+"/payments", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("payment service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code from payment service: %d", resp.StatusCode)
	}

	var pResp paymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&pResp); err != nil {
		return "", fmt.Errorf("failed to decode payment response: %w", err)
	}

	if pResp.Status == "Declined" {
		return "Declined", errors.New("payment declined")
	}

	return pResp.Status, nil
}
