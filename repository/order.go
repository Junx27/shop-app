package repository

import (
	"context"
	"errors"

	"github.com/Junx27/shop-app/entity"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) entity.OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetUserID(id uint) (uint, error) {
	var order entity.Order
	if err := r.db.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("order not found")
		}
		return 0, err
	}
	return order.UserID, nil
}

func (r *OrderRepository) GetManyByUser(ctx context.Context, userID uint, page, limit int) ([]interface{}, error) {
	orders, _, err := r.GetMany(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, len(orders))
	for i, order := range orders {
		result[i] = order
	}
	return result, nil
}

func (r *OrderRepository) GetManyAdmin(ctx context.Context, page, limit int) ([]*entity.Order, int64, error) {
	var orders []*entity.Order
	var total int64
	err := r.db.Model(&entity.Order{}).Count(&total).Offset((page - 1) * limit).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}
func (r *OrderRepository) GetMany(ctx context.Context, userId uint, page, limit int) ([]*entity.Order, int64, error) {
	var orders []*entity.Order
	var total int64
	err := r.db.Model(&entity.Order{}).Where("user_id = ?", userId).Count(&total).Offset((page - 1) * limit).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
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
