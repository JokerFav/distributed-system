package repository

import (
	"order_service/internal/domain"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}
