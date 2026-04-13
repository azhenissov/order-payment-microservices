package service

import (
	"time"
)

type OrderStatusChangeEvent struct {
	OrderID   string
	OldStatus string
	NewStatus string
	UpdatedAt time.Time
	Message   string
}

