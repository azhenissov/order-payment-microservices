# Assignment 1: Clean Architecture Microservices (Order & Payment)

This repository contains the solution for **Assignment 1**, showcasing two microservices (`order-service` and `payment-service`) built in Go following **Clean Architecture** principles.

## Architecture Design & Bounded Contexts

The system is decoupled into two independent bounded contexts with their own data ownership:

1. **Order Service (Port 8080):**
   - Owns the domain of customer orders (`id`, `customer_id`, `item_name`, `amount`, `status`).
   - Uses a dedicated PostgreSQL database (`orders_db`).
   - Communicates synchronously via REST with the Payment Service.
   - Includes custom Idempotency features to prevent duplicate orders.

2. **Payment Service (Port 8081):**
   - Owns the domain of payment transactions (`id`, `order_id`, `transaction_id`, `amount`, `status`).
   - Uses a dedicated PostgreSQL database (`payments_db`).
   - Provides REST APIs to process and retrieve payments. It decides whether to `Authorize` or `Decline` transactions based on business rules.

### Strict Clean Architecture Layers
Both services use strict layers:
- **Domain:** Pure business models, untouched by external tools (`internal/domain`).
- **Use Case:** Encapsulates the specific business rules (e.g., limits, status validation) (`internal/usecase`).
- **Repository:** PostgreSQL implementations fulfilling domain interfaces (`internal/repository`).
- **Transport / HTTP:** Thin handlers parsing requests and calling Use Cases using Gin (`internal/transport/http`).
- **Composition Root:** Manual Dependency Injection wiring the layers in `cmd/main.go`.

### Resilience & Failure Handling

Robust failure handling mechanics have been put in place if the Payment Service becomes unresponsive:
- **HTTP Client Timeout:** The Order Service uses a custom `http.Client` that times out abruptly after **2 seconds**.
- **Service Unavailable (503):** If a timeout occurs (or if the payment service returns another transport error), the Order Service intercepts the error and returns an explicit `503 Service Unavailable` error instead of hanging indefinitely.
- **Fail-Safe Status:** The interrupted Order is immediately updated to `Failed` in the Order database.

### Bonus: Idempotency
To prevent duplicate requests (e.g., from network retries), the Order Service reads the `Idempotency-Key` header. If the exact same key comes in again, it returns the previously generated Order directly without processing another payment or duplicate DB entry.

## Running the Project

### 1. Start the Databases
A `docker-compose.yml` is provided to spin up both distinct databases locally.
```bash
docker-compose up -d
```

### 2. Setup Services
Navigate into each service, install dependencies, and run:

**Payment Service:**
```bash
cd payment-service
go mod tidy
go run cmd/payment-service/main.go
# Server runs on :8081
```

**Order Service:**
```bash
cd order-service
go mod tidy
go run cmd/order-service/main.go
# Server runs on :8080
```
> **Note:** Do not forget to apply the database migrations found inside `migrations/001_initial_schema.sql` of both services against their respective databases before running them.

---

## API Examples (cURL)

### 1. Create a successful Order (Amount < 100,000 cents)
```bash
curl -X POST http://localhost:8080/orders \
-H "Content-Type: application/json" \
-H "Idempotency-Key: request123" \
-d '{"customer_id": "cust-01", "item_name": "Laptop", "amount": 95000}'
```

### 2. Create a Declined Order (Amount > 100,000 cents)
```bash
curl -X POST http://localhost:8080/orders \
-H "Content-Type: application/json" \
-d '{"customer_id": "cust-02", "item_name": "Car", "amount": 5000000}'
```

### 3. Cancel an Order (Must be Pending; if Paid it will fail)
```bash
curl -X PATCH http://localhost:8080/orders/{order_id_here}/cancel
```

### 4. Direct call to Payment Service
```bash
curl -X GET http://localhost:8081/payments/{order_id_here}
```


### Architecture Principles Applied

| Principle | Implementation |
|-----------|-----------------|
| **Layered Architecture** | 4 layers: Transport → UseCase → Domain → Repository |
| **Dependency Inversion** | UseCase depends on interfaces, not implementations |
| **Separation of Concerns** | Each layer has single responsibility |
| **No Framework in Domain** | Domain layer is framework-agnostic |
| **Interface-based Dependencies** | PaymentClient & OrderRepository are interfaces |
| **Manual DI** | Dependencies wired in main.go, no DI framework needed |
| **Bounded Contexts** | Separate databases, no shared code/models |
| **Resilience** | 2-second timeout, 503 on failure, graceful degradation |
| **Idempotency** | Header-based duplicate prevention |
