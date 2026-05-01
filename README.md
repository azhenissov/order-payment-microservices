# Order-Payment Microservices (EDA)

An event-driven microservices system for order processing, payment handling, and notification delivery using asynchronous messaging.

## Overview

This project demonstrates a **choreography-based Event-Driven Architecture (EDA)** where services communicate via events instead of direct calls.

Each service reacts to events and emits new ones, forming a loosely coupled workflow.

## Services

- **Order Service**
  - Creates orders
  - Publishes `OrderCreated` events

- **Payment Service**
  - Consumes `OrderCreated`
  - Processes payment
  - Publishes `PaymentSucceeded` / `PaymentFailed`

- **Notification Service**
  - Consumes payment events
  - Sends notifications

## Architecture

- Event-driven (EDA)
- Choreography-based flow (no central orchestrator)
- Asynchronous communication via RabbitMQ
- Loose coupling between services

### Event Flow

1. Order created → `OrderCreated`
2. Payment Service processes → emits result event  
3. Notification Service reacts → sends notification  

## Tech Stack

- Go (Golang)
- RabbitMQ
- REST APIs
- Docker (optional)

## Getting Started

### Prerequisites

- Go installed
- Docker & Docker Compose (optional)
- RabbitMQ running

### Run

```
docker-compose up -d
```

## Architecture

- Microservices architecture
- Event-driven communication (pub/sub)
- Choreography pattern (no orchestrator)
- Message broker: **RabbitMQ**
- Asynchronous processing
- Eventual consistency

---

## Services

### Order Service
- Accepts HTTP requests to create orders
- Persists order data (if storage is used)
- Publishes `OrderCreated` event

---

### Payment Service
- Subscribes to `OrderCreated`
- Simulates payment processing
- Publishes:
  - `PaymentSucceeded`
  - `PaymentFailed`

---

### Notification Service
- Subscribes to payment events
- Sends notifications (e.g., email/log)
- Implements idempotency protection

---

## Event Flow

```text
Client
  ↓
Order Service → publishes → OrderCreated
  ↓
Payment Service → consumes → OrderCreated
  ↓
Payment Service → publishes → PaymentSucceeded / PaymentFailed
  ↓
Notification Service → consumes → payment events
  ↓
Notification sent
```

