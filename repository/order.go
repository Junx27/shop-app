package repository

import (
	"context"

	"github.com/Junx27/shop-app/entity"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) entity.OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetMany(ctx context.Context, page, limit int) ([]*entity.Order, int64, error) {
	var orders []*entity.Order
	var totalItems int64
	query := r.db.Model(&entity.Order{})

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, totalItems, nil
}

func (r *OrderRepository) CreateOne(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}
