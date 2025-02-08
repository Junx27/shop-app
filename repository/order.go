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
func (r *OrderRepository) GetManyByStatus(ctx context.Context, status string) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := r.db.Where("status = ?", status).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) CreateOne(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) UpdatePayment(ctx context.Context, id uint) (*entity.Order, error) {
	order := &entity.Order{}
	if res := r.db.Model(order).Where("id = ?", id).First(order); res.Error != nil {
		return nil, res.Error
	}

	order.Payment = true
	order.Status = "success"

	if err := r.db.Save(order).Error; err != nil {
		return nil, err
	}

	return order, nil

}
