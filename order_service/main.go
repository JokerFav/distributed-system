package main

import (
	"log"
	"order_service/internal/config"
	"order_service/internal/domain"
	"order_service/internal/handler"
	"order_service/internal/httpapi"
	"order_service/internal/repository"
	"order_service/internal/service"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	db.AutoMigrate(&domain.Order{})

	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)
	e := httpapi.NewRouter(h)

	log.Printf("Order Service is running at %s", cfg.HTTPAddress)
	e.Start(cfg.HTTPAddress)
}
