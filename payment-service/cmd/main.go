package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	"payment-service/internal/api"
	"payment-service/internal/api/middleware"
	"payment-service/internal/repository"
	"payment-service/internal/service"

	desc "github.com/azhenissov/grpc-contracts-go/payment_v1"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load .env file - try multiple paths
	envPaths := []string{".env", "../.env", "../../payment-service/.env"}
	loaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			loaded = true
			log.Printf("Loaded .env from: %s\n", path)
			break
		}
	}
	if !loaded {
		log.Println("Warning: Could not load .env file from any location")
	}

	// 1. Setup Database Connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("Database URL is not in .env")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err := db.Ping()
		if err == nil {
			log.Println("✓ Database connected successfully")
			break
		}
		log.Printf("Failed to connect to database, retrying in 2 seconds... (%d/%d)", i+1, maxRetries)
		if i == maxRetries-1 {
			log.Fatalf("Could not connect to database after %d attempts: %v", maxRetries, err)
		}
		time.Sleep(2 * time.Second)
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		panic("RABBITMQ_URL is not in .env")
	}

	var rabbitConn *amqp.Connection
	var rabbitCh *amqp.Channel

	for i := 0; i < 5; i++ {
		rabbitConn, err = amqp.Dial(rabbitURL)
        if err == nil {
            rabbitCh, err = rabbitConn.Channel()
            if err == nil {
                log.Println("✓ RabbitMQ connected successfully")
                break
            }
        }
        log.Printf("Failed to connect to RabbitMQ, retrying in 3s... (%d/5)", i+1)
        time.Sleep(3 * time.Second)
    }
    if err != nil {
        log.Fatalf("Could not connect to RabbitMQ after retries: %v", err)
    }
    defer rabbitConn.Close()
    defer rabbitCh.Close()

	// 2. Setup Clean Architecture Layers
	paymentRepo := repository.NewPostgresPaymentRepository(db)

	paymentUC := service.NewPaymentUseCase(paymentRepo, rabbitCh) // Передаем rabbitCh в UseCase

	// 3. Configure gRPC Server with Middleware (Interceptors)
	port := os.Getenv("PAYMENT_GRPC_PORT")
	if port == "" {
		panic("PAYMENT_GRPC_PORT is not in .env")
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Add UnaryInterceptor for logging and middleware
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.LoggingUnaryInterceptor),
	)

	paymentHandler := api.NewPaymentGRPCHandler(paymentUC)
	desc.RegisterPaymentAPIServer(grpcServer, paymentHandler)

	log.Printf("✓ gRPC Payment Service started on port %s\n", port)
	log.Println("  Proto contracts: github.com/azhenissov/grpc-contracts-go/payment_v1")
	log.Println("  Documentation: See README.md for details")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
