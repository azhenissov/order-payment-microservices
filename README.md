# Assignment 2: gRPC Microservices with Real-Time Streaming & Event-Driven Architecture

This repository contains the solution for **Assignment 2**, showcasing a comprehensive migration from REST to **gRPC** for inter-service communication, while maintaining **Clean Architecture** principles and backward REST compatibility. The system now features an event-driven architecture with real-time streaming capabilities for order status updates.

## Table of Contents

1. [Overview](#overview)
2. [Key Features](#key-features)
3. [What Changed from Assignment 1](#what-changed-from-assignment-1)
4. [System Architecture](#system-architecture)
5. [Technology Stack](#technology-stack)
6. [Service Details](#service-details)
7. [gRPC & Protocol Buffers](#grpc--protocol-buffers)
8. [Event-Driven Architecture](#event-driven-architecture)
9. [Running the Project](#running-the-project)
10. [Project Structure](#project-structure)
11. [Environment Configuration](#environment-configuration)
12. [API Examples](#api-examples)
13. [Docker & Deployment](#docker--deployment)
14. [Development Guide](#development-guide)

---

## Overview

Assignment 2 represents a significant architectural evolution, transforming the microservices system from REST-based inter-service communication to **gRPC-first** architecture with Protocol Buffers. The system now incorporates:

- **gRPC services** for high-performance inter-service communication
- **Real-time streaming** with event-driven order updates
- **Clean Architecture** patterns with strict separation of concerns
- **Event broker pattern** for asynchronous communication
- **Backward REST API** for client compatibility
- **Docker containerization** for production-ready deployment
- **Environment-based configuration** with `.env` support
- **Middleware interceptors** for request logging and monitoring

**Key Achievements:**
- ✅ Full gRPC migration for service-to-service communication
- ✅ Event-driven architecture with order broker
- ✅ Real-time streaming from database-driven events
- ✅ Proto contracts in dedicated GitHub repositories
- ✅ Clean Architecture preserved (only delivery layer changed)
- ✅ REST API backward compatibility for clients
- ✅ gRPC middleware/interceptor logging
- ✅ Environment-based configuration using godotenv
- ✅ Docker Compose multi-container deployment
- ✅ PostgreSQL with schema migrations

---

## Key Features

### gRPC Communication
- **Protocol Buffers** for efficient serialization
- **Unary RPC** for request-response operations
- **Server-side streaming** for real-time order updates
- **External contract management** via GitHub repositories
- **Code generation** from `.proto` files

### Event-Driven Architecture
- **Order Broker** pattern for event publishing/subscription
- **Order Events**: `OrderCreated`, `OrderStatusChanged`, `OrderCompleted`
- **Real-time streaming** to connected clients
- **Async communication** between services
- **Database-driven event sourcing** capability

### Clean Architecture
- **Domain Layer**: Core business entities and rules
- **Use Case/Service Layer**: Business logic and orchestration
- **Repository Layer**: Data access with PostgreSQL
- **Handler Layer**: Protocol-specific implementations (REST + gRPC)
- **No business logic** in handlers - only protocol adaptation

### Backward Compatibility
- **Dual API support**: REST (HTTP) + gRPC
- **Same port** for REST API (8080)
- **Optional gRPC port** (9091) for streaming
- **Clients can choose** between protocols
- **Gradual migration** path from REST to gRPC

### Deployment
- **Docker containers** for each service
- **Docker Compose** orchestration
- **PostgreSQL databases** per service
- **Network isolation** between services
- **Environment configuration** per environment

### Configuration Management
- **.env files** for environment variables
- **Service-specific configuration** per microservice
- **godotenv** for automatic loading
- **Support for** Docker, local development, CI/CD

---

## What Changed from Assignment 1

### Order Service Changes
| Feature | Assignment 1 | Assignment 2 |
|---------|--------------|-------------|
| **External Communication** | REST (Gin HTTP Client) | gRPC + REST (backward compat) |
| **Order Status Updates** | One-shot HTTP GET | Real-time streaming via gRPC |
| **Payment Integration** | Synchronous REST call | Async gRPC client |
| **Data Format** | JSON | Protocol Buffers (gRPC) |
| **Ports** | 8080 (REST only) | 8080 (REST), 9091 (gRPC) |

### Payment Service Changes
| Feature | Assignment 1 | Assignment 2 |
|---------|--------------|-------------|
| **API Interface** | REST APIs only | gRPC server (Gin removed) |
| **Transport Protocol** | HTTP/JSON | gRPC with Protobuf |
| **Port** | 8081 | 9090 |
| **Request Handling** | Gin handlers | gRPC unary handlers |

### Domain & Use Cases
**No changes** — All business logic remains in the service layer. Only the delivery layer (handler interfaces) changed from HTTP to gRPC.

---

## System Architecture

### High-Level Overview
```
┌──────────────────────────────────────────────────────────┐
│ Client Applications                                      │
│ (Web Browsers, Mobile Apps, Other Services)            │
└────────────────┬──────────────────────┬──────────────────┘
                 │                      │
        ┌────────▼────────┐    ┌────────▼────────┐
        │ REST API        │    │ gRPC Streaming  │
        │ (JSON/HTTP)     │    │ (Binary/Proto)  │
        │ Port: 8080      │    │ Port: 9091      │
        └────────┬────────┘    └────────┬────────┘
                 │                      │
         ┌───────▼──────────────────────▼────────┐
         │                                         │
         │         ORDER SERVICE                  │
         │  ┌──────────────────────────────────┐  │
         │  │ Delivery Layer (Handlers)        │  │
         │  │ ├─ OrderHandler (REST)          │  │
         │  │ └─ OrderGRPCHandler (gRPC)      │  │
         │  ├──────────────────────────────────┤  │
         │  │ Business Logic                   │  │
         │  │ ├─ OrderUseCase                 │  │
         │  │ ├─ OrderBroker (Event Bus)      │  │
         │  │ └─ GRPCPaymentClient (Adapter)  │  │
         │  ├──────────────────────────────────┤  │
         │  │ Data Access                      │  │
         │  │ └─ OrderRepository (PostgreSQL)  │  │
         │  └──────────────────────────────────┘  │
         │                                         │
         └────────┬────────────────┬───────────────┘
                  │                │
        ┌─────────▼─────────┐      │ gRPC Call
        │ Order Events      │      │ (Unary)
        │ - OrderCreated    │      │
        │ - StatusChanged   │      │
        │ - OrderCompleted  │      │
        │ (Published to     │      │
        │  Streaming API)   │      │
        └───────────────────┘      │
                                   │
                        ┌──────────▼────────────────┐
                        │                           │
                        │  PAYMENT SERVICE          │
                        │ ┌──────────────────────┐ │
                        │ │ Delivery Layer       │ │
                        │ │ └─ PaymentGRPC       │ │
                        │ │    Handler           │ │
                        │ ├──────────────────────┤ │
                        │ │ Business Logic       │ │
                        │ │ └─ PaymentUseCase    │ │
                        │ ├──────────────────────┤ │
                        │ │ Data Access          │ │
                        │ │ └─ PaymentRepository │ │
                        │ └──────────────────────┘ │
                        │ gRPC Server              │
                        │ Port: 9090               │
                        └──────────────────────────┘
```

### Data Flow: Creating an Order
```
1. Client sends REST POST /orders
                 ↓
2. OrderHandler receives & validates
                 ↓
3. OrderUseCase processes:
   - Creates order in database
   - Calls PaymentGRPCClient (gRPC Unary)
   - Publishes OrderCreated event to Broker
                 ↓
4. PaymentService (via gRPC):
   - Processes payment
   - Returns status
                 ↓
5. OrderBroker emits event:
   - All subscribers notified (e.g., streaming clients)
   - OrderStatusChanged events go to listeners
                 ↓
6. Response returned with order ID + status
```

### Event Flow: Real-Time Streaming
```
Client A: gRPC Stream Subscribe
Client B: gRPC Stream Subscribe
                   ↓
         ┌─────────┴──────────┐
         │ OrderGRPCHandler   │
         │ (Server Streaming) │
         └─────────┬──────────┘
                   ↓
          OrderBroker listens to events
                   ↓
        ┌──────────┴───────────┐
        │                      │
Order updates flow to    All subscribed clients
connected clients        receive real-time updates
```

---

## Technology Stack

### Languages & Frameworks
- **Language**: Go 1.25.5
- **REST Framework**: Gin Gonic (HTTP server)
- **gRPC**: google.golang.org/grpc (RPC framework)
- **Serialization**: Protocol Buffers 3.x

### Database & Persistence
- **DBMS**: PostgreSQL 15 (Alpine)
- **Driver**: lib/pq (pure Go PostgreSQL driver)
- **Migrations**: SQL files in `/migrations`
- **per-service databases**: Strict data isolation

### Infrastructure & Deployment
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Environment Management**: godotenv (.env files)
- **CI/CD**: GitHub Actions (configured)

### External Dependencies
- **Protocol Contracts**: https://github.com/azhenissov/grpc-contracts-go
  - `order_v1` - Order Service contracts
  - `payment_v1` - Payment Service contracts

---

## Service Details

### Order Service (Port 8080/9091)
**Responsibilities:**
- Manage orders from creation to completion
- Communicate with Payment Service via gRPC
- Stream order status updates in real-time
- Maintain order history and state

**Exposed APIs:**
- **REST** (Port 8080):
  - `POST /orders` - Create new order
  - `GET /orders/:id` - Get order details
  - `GET /orders` - List all orders
  
- **gRPC** (Port 9091):
  - `CreateOrder` - Unary RPC for order creation
  - `GetOrder` - Unary RPC for order retrieval
  - `StreamOrders` - Server-side streaming for real-time updates

**Database:**
- `orders_db` on PostgreSQL
- Tables: `orders`, `order_items`, `order_status_history`, `events`

**Dependencies:**
- PostgreSQL (orders_db)
- Payment Service (gRPC client)

### Payment Service (Port 9090)
**Responsibilities:**
- Process payments for orders
- Validate payment information
- Track payment history
- Pure gRPC server (no REST API)

**Exposed APIs:**
- **gRPC** (Port 9090):
  - `ProcessPayment` - Unary RPC for payment processing
  - `GetPayment` - Unary RPC for payment lookup
  - `GetPaymentHistory` - Unary RPC for payment records

**Database:**
- `payments_db` on PostgreSQL
- Tables: `payments`, `payment_methods`, `transactions`

**Dependencies:**
- PostgreSQL (payments_db)

---

## gRPC & Protocol Buffers

### Proto Files Location
```
Proto contracts are maintained externally in GitHub:
https://github.com/azhenissov/grpc-contracts-go

Generated Go code is in each service:
├── order-service/pkg/pb/order_v1/
│   ├── order.pb.go          (Messages)
│   └── order_grpc.pb.go     (Service definitions)
└── payment-service/pkg/pb/payment_v1/
    ├── payment.pb.go        (Messages)
    └── payment_grpc.pb.go   (Service definitions)
```

### Service Definitions
```protobuf
// Order Service
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
  rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
  rpc StreamOrders(Empty) returns (stream OrderUpdate);
}

// Payment Service
service PaymentAPI {
  rpc ProcessPayment(PaymentRequest) returns (PaymentResponse);
  rpc GetPayment(GetPaymentRequest) returns (GetPaymentResponse);
}
```

### Message Types
```protobuf
message Order {
  string id = 1;
  string customer_id = 2;
  repeated OrderItem items = 3;
  string status = 4;
  double total_amount = 5;
  google.protobuf.Timestamp created_at = 6;
}

message Payment {
  string id = 1;
  string order_id = 2;
  double amount = 3;
  string status = 4;
  string payment_method = 5;
}
```

---

## Event-Driven Architecture

### Order Broker Pattern
The **OrderBroker** is the central event bus for the Order Service:

```go
type OrderBroker interface {
  Subscribe(eventType string) <-chan interface{}
  Publish(event interface{})
  Unsubscribe(eventType string, ch <-chan interface{})
}
```

### Event Types
```go
// OrderCreated - emitted when new order is created
type OrderCreated struct {
  OrderID    string
  CustomerID string
  Amount     float64
  Timestamp  time.Time
}

// OrderStatusChanged - emitted when order status updates
type OrderStatusChanged struct {
  OrderID    string
  OldStatus  string
  NewStatus  string
  Timestamp  time.Time
}

// OrderCompleted - emitted when order is finalized
type OrderCompleted struct {
  OrderID    string
  Timestamp  time.Time
}
```

### Event Flow
```
1. Use Case creates order → Publishes OrderCreated
             ↓
2. Broker receives event → Notifies all subscribers
             ↓
3. gRPC Streaming Handler:
   - Receives event
   - Sends OrderUpdate to connected clients
   - Multiple clients see real-time updates
             ↓
4. Database layer (optional):
   - Events saved to events table
   - Enables event sourcing
   - Audit trail maintained
```

### Subscribing to Events
```go
// In gRPC Handler
broker := service.GetOrderBroker()
eventChan := broker.Subscribe("OrderStatusChanged")

// Stream to client
for event := range eventChan {
  update := convertToProto(event)
  stream.Send(update)
}
```

---

## Running the Project

### Prerequisites
- Go 1.25.5+
- Docker & Docker Compose
- PostgreSQL client tools (optional)

### Quick Start with Docker

```bash
# Clone the repository
git clone https://github.com/yourusername/AP2_1.git
cd AP2_1

# Start all services with Docker Compose
docker-compose up --build

# Services will be available:
# - Order Service REST: http://localhost:8080
# - Order Service gRPC: localhost:9091
# - Payment Service gRPC: localhost:9090
```

### Local Development

```bash
# Navigate to order service
cd order-service

# Load environment variables (.env)
# Ensure PostgreSQL is running on localhost:5434

# Run the service
go run ./cmd/main.go

# In another terminal, run payment service
cd ../payment-service
go run ./cmd/main.go
```

### Environment Setup

Each service has its own `.env` file:

**order-service/.env:**
```env
ORDER_REST_PORT=:8080
ORDER_GRPC_PORT=:9091
DATABASE_URL=postgres://order_user:order_password@localhost:5434/orders_db?sslmode=disable
PAYMENT_GRPC_ADDRESS=localhost:9090
```

**payment-service/.env:**
```env
PAYMENT_GRPC_PORT=:9090
DATABASE_URL=postgres://payment_user:payment_password@localhost:5433/payments_db?sslmode=disable
```

### Database Setup

The Docker Compose automatically:
1. Creates PostgreSQL containers
2. Creates databases
3. Applies migrations from `/migrations` folders

For manual setup:
```bash
# Create databases
createdb -h localhost -U order_user orders_db
createdb -h localhost -U payment_user payments_db

# Apply migrations
psql -h localhost -U order_user -d orders_db < order-service/migrations/001_initial_schema.sql
psql -h localhost -U payment_user -d payments_db < payment-service/migrations/001_initial_schema.sql
```

---

## Project Structure

```
AP2_1/
├── README.md
├── docker-compose.yml          # Docker Compose configuration
├── .github/                     # GitHub Actions workflows
│   └── workflows/
│       └── ci-cd.yml
│
├── order-service/
│   ├── .env                     # Environment configuration
│   ├── go.mod                   # Module dependencies
│   ├── go.sum                   # Dependency checksums
│   ├── Dockerfile               # Docker image definition
│   ├── cmd/
│   │   └── main.go             # Application entry point
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler.go      # REST HTTP handlers
│   │   │   ├── grpc_order_handler.go    # gRPC server implementation
│   │   │   └── grpc_payment_client.go   # gRPC client to Payment Service
│   │   ├── config/
│   │   │   └── config.go       # Configuration management
│   │   ├── domain/
│   │   │   └── order.go        # Order entity and business rules
│   │   ├── repository/
│   │   │   └── order_postgres.go  # PostgreSQL implementation
│   │   └── service/
│   │       ├── order_usecase.go    # Business logic orchestration
│   │       ├── order_event.go      # Event definitions
│   │       └── order-broker.go     # Event broker/bus
│   ├── pkg/pb/
│   │   └── order_v1/           # Generated Protocol Buffer code
│   │       ├── order.pb.go
│   │       └── order_grpc.pb.go
│   └── migrations/
│       └── 001_initial_schema.sql
│
├── payment-service/
│   ├── .env                     # Environment configuration
│   ├── go.mod                   # Module dependencies
│   ├── go.sum                   # Dependency checksums
│   ├── Dockerfile               # Docker image definition
│   ├── cmd/
│   │   └── main.go             # Application entry point
│   ├── internal/
│   │   ├── api/
│   │   │   ├── grpc_handler.go # gRPC server implementation
│   │   │   └── middleware/
│   │   │       └── interceptor.go  # gRPC interceptors for logging
│   │   ├── config/
│   │   │   └── config.go       # Configuration management
│   │   ├── domain/
│   │   │   └── payment.go      # Payment entity and business rules
│   │   ├── repository/
│   │   │   └── payment_postgres.go # PostgreSQL implementation
│   │   └── service/
│   │       └── payment_usecase.go  # Business logic
│   ├── pkg/pb/
│   │   └── payment_v1/         # Generated Protocol Buffer code
│   │       ├── payment.pb.go
│   │       └── payment_grpc.pb.go
│   └── migrations/
│       └── 001_initial_schema.sql
│
└── diagrams/                    # Architecture diagrams
```

---

## Environment Configuration

### Multi-Service Configuration
Each microservice manages its own configuration independently:

```yaml
Order Service:
  - REST API port: 8080
  - gRPC port: 9091
  - Database: orders_db (port 5434)
  - Config file: order-service/.env

Payment Service:
  - gRPC port: 9090
  - Database: payments_db (port 5433)
  - Config file: payment-service/.env
  - No REST API (pure gRPC)
```

### Loading Configuration
The `.env` files are automatically loaded using `godotenv`:

```go
// In main.go
if err := godotenv.Load("../.env"); err != nil {
  log.Printf("Warning: Could not load .env file: %v\n", err)
}

// Environment variables are now available
dsn := os.Getenv("DATABASE_URL")
port := os.Getenv("ORDER_REST_PORT")
grpcAddr := os.Getenv("PAYMENT_GRPC_ADDRESS")
```

---

## API Examples

### REST API (Order Service)

**Create Order:**
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cust123",
    "items": [
      {"product_id": "prod1", "quantity": 2, "price": 29.99}
    ]
  }'
```

Response:
```json
{
  "id": "order-uuid",
  "customer_id": "cust123",
  "status": "pending",
  "total_amount": 59.98,
  "created_at": "2026-04-13T17:11:47Z"
}
```

**Get Order:**
```bash
curl http://localhost:8080/orders/order-uuid
```

**List Orders:**
```bash
curl http://localhost:8080/orders
```

### gRPC API (Order Service - Streaming Example)

Using `grpcurl`:

```bash
# Stream order updates
grpcurl -plaintext localhost:9091 order_v1.OrderService.StreamOrders

# Create order (unary)
grpcurl -plaintext \
  -d '{"customer_id":"cust123","items":[{"product_id":"prod1","quantity":2}]}' \
  localhost:9091 order_v1.OrderService.CreateOrder
```

### gRPC Communication Between Services

Order Service → Payment Service (internal):

```go
// OrderUseCase calls Payment Service via gRPC
paymentReq := &paymentDesc.PaymentRequest{
  OrderID: order.ID,
  Amount:  order.TotalAmount,
}

paymentResp, err := o.paymentClient.ProcessPayment(ctx, paymentReq)
if err != nil {
  return err
}

// Response received, continue with business logic
```

---

## Docker & Deployment

### Docker Compose Services
```yaml
Services:
  - order_db: PostgreSQL 15 (port 5434)
  - payment_db: PostgreSQL 15 (port 5433)
  - order-service: Go application (ports 8080, 9091)
  - payment-service: Go application (port 9090)

Volumes:
  - order_data: Persists order database
  - payment_data: Persists payment database
```

### Building Images

```bash
# Build individual images
docker build -t order-service:latest ./order-service
docker build -t payment-service:latest ./payment-service

# Build with Docker Compose (automatic)
docker-compose build

# Run services
docker-compose up
```

### Production Considerations
- Use `.env.prod` for production variables
- Enable gRPC TLS/SSL certificates
- Configure resource limits in Docker
- Use secrets management (not plain .env)
- Enable monitoring and logging
- Implement service health checks

---

## Development Guide

### Adding New gRPC Methods

1. **Define proto message** in external contracts repo:
```protobuf
service OrderService {
  rpc NewMethod(NewRequest) returns (NewResponse);
}
```

2. **Regenerate code:**
```bash
cd grpc-contracts-go
protoc --go_out=. --go-grpc_out=. order_v1/order.proto
```

3. **Update handler in order-service:**
```go
func (h *OrderGRPCHandler) NewMethod(ctx context.Context, req *orderDesc.NewRequest) (*orderDesc.NewResponse, error) {
  // Implementation
}
```

4. **Update use case:**
```go
func (u *OrderUseCase) NewBusinessLogic(ctx context.Context, input interface{}) error {
  // Business logic
}
```

### Adding New Events

1. **Define event structure:**
```go
type NewEvent struct {
  OrderID   string
  Timestamp time.Time
}
```

2. **Publish from use case:**
```go
event := &NewEvent{OrderID: order.ID, Timestamp: time.Now()}
u.broker.Publish(event)
```

3. **Subscribe in handler/service:**
```go
eventChan := u.broker.Subscribe("NewEvent")
for event := range eventChan {
  // Handle event
}
```

### Testing

```bash
# Unit tests
go test ./... -v

# Integration tests with Docker
docker-compose up -d
go test ./... -v -tags=integration
docker-compose down
```

### Database Migrations

Create new migration:
```sql
-- migrations/002_new_feature.sql
BEGIN;

ALTER TABLE orders ADD COLUMN new_field VARCHAR(255);

COMMIT;
```

Apply with Docker (automatic) or manually:
```bash
psql -h localhost -U user -d database < migrations/002_new_feature.sql
```

---

### Clean Architecture Layers (Unchanged Except Delivery)

Both services maintain strict Clean Architecture:

**1. Domain Layer** (`internal/domain/`)
- Pure business models: `Order`, `Payment`
- Domain interfaces: `OrderRepository`, `PaymentRepository`, `PaymentClient`
- **No framework dependencies**

**2. Service/Use Case Layer** (`internal/service/`)
- Business rules and orchestration
- `OrderUseCase`: Create order, get order, cancel order, revenue calculation
- `PaymentUseCase`: Process payment, get payment status
- **Entirely unchanged from Assignment 1**

**3. Repository Layer** (`internal/repository/`)
- PostgreSQL data access implementations
- Fulfills domain interfaces
- **Unchanged from Assignment 1**

**4. Delivery/API Layer** (`internal/api/`)
- **gRPC Handlers** (NEW):
  - `grpc_order_handler.go`: Implements `OrderServiceServer`
  - `grpc_payment_client.go`: Implements `domain.PaymentClient` using gRPC
- **REST Handlers** (PRESERVED):
  - `handler.go`: Gin HTTP handlers for backward compatibility
- Thin adapters converting transport format to domain models

**5. Composition Root** (`cmd/main.go`)
- Dependency injection and wiring
- Creates both gRPC and REST servers
- Configures environment variables

---

## Contract Repositories

### Proto Definitions
**Repository:** https://github.com/azhenissov/grpc-contracts-proto

Contains Protocol Buffer service definitions:
- `payment/v1/payment.proto`: PaymentAPI service definition
- `order/v1/order.proto`: OrderService service definition with streaming

### Generated Go Code
**Repository:** https://github.com/azhenissov/grpc-contracts-go

Generated gRPC stubs and message types using `protoc`:
- `payment_v1/`: Generated PaymentAPI client/server interfaces
- `order_v1/`: Generated OrderService client/server interfaces

### Local Proto Files (Reference)
For documentation purposes, proto files are also stored locally:
- `payment-service/proto/payment/v1/payment.proto`
- `order-service/proto/order/v1/order.proto`

---

## Service Details

### Payment Service (gRPC)

**Port:** `9090` (via `PAYMENT_GRPC_PORT` env var)

**gRPC Service:** `PaymentAPI`

**Methods:**
- `ProcessPayment(ProcessPaymentRequest) → ProcessPaymentResponse`
  - Input: `order_id`, `amount`
  - Output: `transaction_id`, `status` (AUTHORIZED or DECLINED)
  - Business Rule: Amounts < 100,000 cents are authorized; ≥ 100,000 are declined

- `GetPaymentStatus(GetPaymentStatusRequest) → GetPaymentStatusResponse`
  - Input: `order_id`
  - Output: `status`, `transaction_id`

**Handler Location:** `internal/api/grpc_handler.go`

**Import Path:** 
```go
import (
    desc "github.com/azhenissov/grpc-contracts-go/payment_v1"
)
```

---

### Order Service (Hybrid REST + gRPC)

**REST Port:** `8080` (via `ORDER_REST_PORT` env var)  
**gRPC Port:** `9091` (via `ORDER_GRPC_PORT` env var)  
**Payment gRPC Address:** `PAYMENT_GRPC_ADDRESS` (default: `localhost:9090`)

**gRPC Service:** `OrderService`

**Methods:**
- `CreateOrder(CreateOrderRequest) → CreateOrderResponse`
  - Input: `customer_id`, `item_name`, `amount`
  - Output: `order_id`, `status`, `created_at`
  - Calls Payment Service gRPC for payment processing
  - **Streaming Events:** Order status changes trigger DB events

- `GetOrder(GetOrderRequest) → Order`
  - Input: `order_id`
  - Output: `Order` message with full details

- `SubscribeToOrderUpdates(OrderSubscriptionRequest) → stream OrderStatusUpdate` *(Server-side streaming)*
  - Input: `order_id`
  - Output: Continuous stream of `OrderStatusUpdate` messages
  - **Real-Time:** Database changes trigger channel events sent to subscribers

**REST API (Backward Compatibility):**
- `POST /orders`: Create order (Gin handler)
- `GET /orders/:id`: Get order
- `PATCH /orders/:id/cancel`: Cancel order
- `GET /revenue`: Get total revenue

**Handler Locations:**
- gRPC: `internal/api/grpc_order_handler.go`
- REST: `internal/api/handler.go`
- Payment Client Adapter: `internal/api/grpc_payment_client.go`

**Import Paths:**
```go
import (
    desc "github.com/azhenissov/grpc-contracts-go/order_v1"
    paymentDesc "github.com/azhenissov/grpc-contracts-go/payment_v1"
)
```

---

## Running the Project

### Prerequisites
- Go 1.20+
- PostgreSQL (or Docker/Docker Compose)
- `protoc` (optional, only if regenerating from proto files)

### 1. Start Databases

```bash
docker-compose up -d
```

This starts:
- `orders_db`: PostgreSQL on port 5433 (Order Service)
- `payments_db`: PostgreSQL on port 5434 (Payment Service)

### 2. Apply Migrations

**Payment Service:**
```bash
psql -h localhost -p 5434 -U postgres -d payments_db < payment-service/migrations/001_initial_schema.sql
```

**Order Service:**
```bash
psql -h localhost -p 5433 -U postgres -d orders_db < order-service/migrations/001_initial_schema.sql
```

### 3. Run Services

**Terminal 1 - Payment Service:**
```bash
cd payment-service
go mod tidy
PAYMENT_GRPC_PORT=:9090 go run cmd/main.go
```

**Terminal 2 - Order Service:**
```bash
cd order-service
go mod tidy
ORDER_REST_PORT=:8080 ORDER_GRPC_PORT=:9091 PAYMENT_GRPC_ADDRESS=localhost:9090 go run cmd/main.go
```

**Expected Output:**
```
Payment Service: gRPC server listening on :9090
Order Service: REST server listening on :8080
Order Service: gRPC server listening on :9091
```

---

## API Examples

### REST API (Backward Compatibility)

#### Create Order via REST
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cust-001",
    "item_name": "Laptop",
    "amount": 95000
  }'
```

**Response:**
```json
{
  "id": "ord-uuid-001",
  "customer_id": "cust-001",
  "item_name": "Laptop",
  "amount": 95000,
  "status": "PAID",
  "created_at": "2025-01-15T10:30:00Z"
}
```

#### Get Order via REST
```bash
curl http://localhost:8080/orders/ord-uuid-001
```

#### Get Revenue via REST
```bash
curl http://localhost:8080/revenue
```

### gRPC API (Using grpcurl)

#### Install grpcurl
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

#### Create Order via gRPC
```bash
grpcurl -plaintext \
  -d '{"customer_id":"cust-002","item_name":"Mouse","amount":50000}' \
  github.com/azhenissov/grpc-contracts-go.order_v1.OrderService/CreateOrder \
  localhost:9091
```

#### Get Order via gRPC
```bash
grpcurl -plaintext \
  -d '{"order_id":"ord-uuid-001"}' \
  localhost:9091 \
  github.com.azhenissov.grpc-contracts-go.order_v1.OrderService/GetOrder
```

#### Subscribe to Order Updates (Streaming)
```bash
grpcurl -plaintext \
  -d '{"order_id":"ord-uuid-001"}' \
  localhost:9091 \
  github.com.azhenissov.grpc-contracts-go.order_v1.OrderService/SubscribeToOrderUpdates
```

**Streaming Output (real-time as status changes):**
```json
{
  "order_id": "ord-uuid-001",
  "old_status": "PENDING",
  "new_status": "AUTHORIZED",
  "updated_at": "2025-01-15T10:30:15Z"
}
{
  "order_id": "ord-uuid-001",
  "old_status": "AUTHORIZED",
  "new_status": "PAID",
  "updated_at": "2025-01-15T10:30:20Z"
}
```

---

## Streaming Implementation

### Real-Time Database Events

The streaming feature is **NOT simulated** — it's tied to real database state changes:

**Order Service Streaming Architecture:**

1. **Event Generation:** When an order status changes in the database (via use case layer), an `OrderStatusChangeEvent` is published to an in-memory channel:
   ```go
   orderNotifications <- &domain.OrderStatusChangeEvent{
       OrderID:   orderID,
       OldStatus: oldStatus,
       NewStatus: newStatus,
       UpdatedAt: time.Now(),
   }
   ```

2. **Subscription Handling:** Each gRPC client subscription creates a listener:
   ```go
   // In SubscribeToOrderUpdates handler
   for event := range clientUpdates {
       stream.Send(&desc.OrderStatusUpdate{
           OrderId:    event.OrderID,
           OldStatus:  event.OldStatus,
           NewStatus:  event.NewStatus,
           UpdatedAt:  timestamppb.New(event.UpdatedAt),
       })
   }
   ```

3. **Database-Driven Flow:**
   - User creates/updates order → Database updated
   - Use case publishes event to channel
   - All active subscriptions receive the event
   - Clients see real-time updates without polling



## Project Structure

```
AP2_1/
├── README.md                    # This file
├── docker-compose.yml           # Database setup for local development
│
├── order-service/
│   ├── cmd/
│   │   └── main.go             # gRPC + REST server entry point
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler.go      # REST handlers (Gin) - backward compat
│   │   │   ├── grpc_order_handler.go    # gRPC OrderServiceServer impl
│   │   │   └── grpc_payment_client.go   # gRPC PaymentClient adapter
│   │   ├── service/
│   │   │   ├── order_usecase.go        # Business logic (unchanged)
│   │   │   └── order_event.go          # Event struct for streaming
│   │   ├── repository/
│   │   │   └── order_postgres.go       # PostgreSQL data access
│   │   ├── domain/
│   │   │   └── order.go                # Domain models
│   │   └── config/
│   │       └── config.go               # Configuration loading
│   ├── proto/
│   │   └── order/v1/
│   │       └── order.proto     # Proto definition (reference)
│   ├── migrations/
│   │   └── 001_initial_schema.sql
│   └── go.mod
│
├── payment-service/
│   ├── cmd/
│   │   └── main.go             # gRPC server entry point
│   ├── internal/
│   │   ├── api/
│   │   │   ├── grpc_handler.go  # gRPC PaymentAPIServer impl
│   │   │   └── middleware/
│   │   │       └── interceptor.go       # LoggingUnaryInterceptor (bonus)
│   │   ├── service/
│   │   │   └── payment_usecase.go       # Business logic (unchanged)
│   │   ├── repository/
│   │   │   └── payment_postgres.go      # PostgreSQL data access
│   │   ├── domain/
│   │   │   └── payment.go               # Domain models
│   │   └── config/
│   │       └── config.go               # Configuration loading
│   ├── proto/
│   │   └── payment/v1/
│   │       └── payment.proto   # Proto definition (reference)
│   ├── migrations/
│   │   └── 001_initial_schema.sql
│   └── go.mod
│
└── diagrams/                    # Architecture diagrams (if any)
```

---

## Key Implementation Details

### Environment Variables

**Payment Service:**
```bash
PAYMENT_GRPC_PORT=:9090
DATABASE_URL=postgres://user:pass@localhost:5434/payments_db
```

**Order Service:**
```bash
ORDER_REST_PORT=:8080
ORDER_GRPC_PORT=:9091
PAYMENT_GRPC_ADDRESS=localhost:9090
DATABASE_URL=postgres://user:pass@localhost:5433/orders_db
```

### gRPC Server Setup

**Payment Service (`cmd/main.go`):**
```go
grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(middleware.LoggingUnaryInterceptor),
)
paymentHandler := api.NewPaymentGRPCHandler(paymentUC)
desc.RegisterPaymentAPIServer(grpcServer, paymentHandler)
grpcServer.Serve(listener)
```

**Order Service (`cmd/main.go`):**
```go
// gRPC Client for Payment Service
paymentConn, _ := grpc.Dial(paymentGRPCAddr, grpc.WithInsecure())
paymentClient := paymentDesc.NewPaymentAPIClient(paymentConn)

// gRPC Server for Order Service
grpcServer := grpc.NewServer()
orderHandler := api.NewOrderGRPCHandler(orderUC, orderNotifications)
desc.RegisterOrderServiceServer(grpcServer, orderHandler)
grpcServer.Serve(listener)
```

### Streaming with Channels

```go
// In grpc_order_handler.go
func (h *OrderGRPCHandler) SubscribeToOrderUpdates(
    req *desc.OrderSubscriptionRequest,
    stream desc.OrderService_SubscribeToOrderUpdatesServer,
) error {
    clientUpdates := make(chan *domain.OrderStatusChangeEvent)
    h.subscriptions[req.OrderId] = append(h.subscriptions[req.OrderId], clientUpdates)
    
    for event := range clientUpdates {
        stream.Send(&desc.OrderStatusUpdate{
            OrderId:   event.OrderID,
            NewStatus: event.NewStatus,
            // ...
        })
    }
}

// When order status changes in use case:
h.orderNotifications <- &domain.OrderStatusChangeEvent{
    OrderID:   orderID,
    OldStatus: oldStatus,
    NewStatus: newStatus,
}
```

---

## Testing & Validation

### Integration Testing
To verify the gRPC migration works correctly:

1. **Payment Service Directly:**
   ```bash
   grpcurl -plaintext -d '{"order_id":"test-001","amount":50000}' \
     localhost:9090 github.com.azhenissov.grpc-contracts-go.payment_v1.PaymentAPI/ProcessPayment
   ```

2. **Order Service gRPC:**
   ```bash
   grpcurl -plaintext -d '{"customer_id":"cust-003","item_name":"Keyboard","amount":30000}' \
     localhost:9091 github.com.azhenissov.grpc-contracts-go.order_v1.OrderService/CreateOrder
   ```

3. **Streaming Verification:**
   - Create an order via gRPC
   - Start a streaming subscription in another terminal
   - Verify order status changes appear in the stream (in real-time as DB updates)

---

## Summary

This assignment demonstrates a successful migration from REST to gRPC while maintaining Clean Architecture principles. The solution showcases:

- **Protocol Buffer contracts** in dedicated GitHub repositories
- **gRPC server/client** implementations with proper error handling
- **Real-time streaming** tied to database events
- **Request logging** middleware as a bonus feature
- **Backward REST compatibility** for seamless client migration
- **Environment-based configuration** with no hardcoded values
- **Clean separation** of concerns maintained across layers


