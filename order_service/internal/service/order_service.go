package service

import (
	"order_service/internal/domain"
	"order_service/internal/repository"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(userID uint, itemName string, amount float64) (*domain.Order, error) {
	order := &domain.Order{
		UserID:   userID,
		ItemName: itemName,
		Amount:   amount,
		Status:   "PENDING",
	}
	if err := s.repo.Create(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetOrdersByUser(userID uint) ([]domain.Order, error) {
	return s.repo.FindByUserID(userID)
}
