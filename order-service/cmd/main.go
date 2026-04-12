package main

import (
	"database/sql"
	"log"
	"os"

	"order-service/internal/api"
	"order-service/internal/repository"
	"order-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost user=order_user password=order_password dbname=orders_db port=5434 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	paymentSvcURL := os.Getenv("PAYMENT_SERVICE_URL")
	if paymentSvcURL == "" {
		paymentSvcURL = "http://localhost:8081"
	}

	orderRepo := repository.NewPostgresOrderRepository(db)
	paymentClient := api.NewHTTPPaymentClient(paymentSvcURL)
	orderUC := service.NewOrderUseCase(orderRepo, paymentClient)

	router := gin.Default()

	api.NewOrderHandler(router, orderUC)

	log.Println("Order Service running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
