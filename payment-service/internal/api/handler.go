package api

import (
	"net/http"
	"payment-service/internal/domain"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	PaymentUseCase domain.PaymentUseCase
}

func NewPaymentHandler(router *gin.Engine, us domain.PaymentUseCase) {
	handler := &PaymentHandler{
		PaymentUseCase: us,
	}

	router.POST("/payments", handler.ProcessPayment)
	router.GET("/payments/:order_id", handler.GetPaymentStatus)
}

type processPaymentRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	Amount  int64  `json:"amount" binding:"required"`
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req processPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.PaymentUseCase.ProcessPayment(c.Request.Context(), req.OrderID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Logic: Return "Authorized" + unique transaction_id.
	// The assignment specifies these in response:
	c.JSON(http.StatusOK, gin.H{
		"status":         payment.Status,
		"transaction_id": payment.TransactionID,
		"order_id":       payment.OrderID,
	})
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	orderID := c.Param("order_id")

	payment, err := h.PaymentUseCase.GetPaymentStatus(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}
