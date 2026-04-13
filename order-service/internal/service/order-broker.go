package service

import "sync"

type OrderStatus struct {
	OrderID string
	Status  string
}

type OrderBroker struct {
	mu          sync.RWMutex
	subscribers map[string][]chan OrderStatus 
}

func NewOrderBroker() *OrderBroker {
	return &OrderBroker{
		subscribers: make(map[string][]chan OrderStatus),
	}
}

func (b *OrderBroker) Subscribe(orderID string) chan OrderStatus {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan OrderStatus, 1) 
	b.subscribers[orderID] = append(b.subscribers[orderID], ch)
	return ch
}

func (b *OrderBroker) Unsubscribe(orderID string, ch chan OrderStatus) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subs := b.subscribers[orderID]
	for i, sub := range subs {
		if sub == ch {
			// Удаляем канал из слайса
			b.subscribers[orderID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}
}

func (b *OrderBroker) Publish(orderID string, status string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if subs, ok := b.subscribers[orderID]; ok {
		for _, ch := range subs {
			select {
			case ch <- OrderStatus{OrderID: orderID, Status: status}:
			default:
			}
		}
	}
}