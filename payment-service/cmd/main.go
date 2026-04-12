package main

import (
	"database/sql"
	"log"

	"payment-service/internal/api"
	"payment-service/internal/repository"
	"payment-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Connect to PostgreSQL DB
	dsn := "host=localhost user=payment_user password=payment_password dbname=payments_db port=5433 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 2. Setup Clean Architecture Layers
	paymentRepo := repository.NewPostgresPaymentRepository(db)
	paymentUC := service.NewPaymentUseCase(paymentRepo)

	router := gin.Default()

	// 3. Register HTTP handlers
	api.NewPaymentHandler(router, paymentUC)

	log.Println("Payment Service running on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
