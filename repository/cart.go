package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Junx27/shop-app/entity"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) entity.CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetUserID(id uint) (uint, error) {
	var cart entity.Cart
	if err := r.db.First(&cart, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("cart not found")
		}
		return 0, err
	}
	return cart.UserID, nil
}

func (r *CartRepository) GetManyByUser(ctx context.Context, userID uint, page, limit int) ([]interface{}, error) {
	carts, _, err := r.GetMany(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, len(carts))
	for i, cart := range carts {
		result[i] = cart
	}
	return result, nil
}

func (r *CartRepository) GetManyAdmin(ctx context.Context, page, limit int) ([]*entity.Cart, int64, error) {
	var carts []*entity.Cart
	var total int64
	err := r.db.Model(&entity.Cart{}).Count(&total).Offset((page - 1) * limit).Limit(limit).Find(&carts).Error
	if err != nil {
		return nil, 0, err
	}
	return carts, total, nil
}
func (r *CartRepository) GetMany(ctx context.Context, userId uint, page, limit int) ([]*entity.Cart, int64, error) {
	var carts []*entity.Cart
	var total int64
	err := r.db.Model(&entity.Cart{}).Where("user_id = ?", userId).Count(&total).Offset((page - 1) * limit).Limit(limit).Find(&carts).Error
	if err != nil {
		return nil, 0, err
	}
	return carts, total, nil
}

func (r *CartRepository) GetOne(ctx context.Context, id uint) (*entity.Cart, error) {
	cart := &entity.Cart{}
	if res := r.db.Model(cart).Where("id = ?", id).First(cart); res.Error != nil {
		return nil, res.Error
	}

	return cart, nil
}

func (r *CartRepository) FindByUserAndMenuAndStatus(ctx context.Context, userID uint, menuID uint, status string) (*entity.Cart, error) {
	cart := &entity.Cart{}
	if err := r.db.Where("user_id = ? AND menu_id = ? AND status = ?", userID, menuID, status).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return cart, nil
}

func (r *CartRepository) GetManyByUserAndStatus(ctx context.Context, userID uint, status string) ([]*entity.Cart, error) {
	var carts []*entity.Cart
	if err := r.db.Where("user_id = ? AND status = ?", userID, status).Find(&carts).Error; err != nil {
		return nil, err
	}
	return carts, nil
}

func (r *CartRepository) UpdateOrderIDByStatus(ctx context.Context, userID uint, orderID uint) error {
	carts, err := r.GetManyByUserAndStatus(ctx, userID, "pending")
	if err != nil {
		return err
	}
	for _, cart := range carts {
		cart.OrderID = &orderID
		cart.Status = "checkout"
		_, err := r.UpdateOne(ctx, cart)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CartRepository) UpdateOne(ctx context.Context, cart *entity.Cart) (*entity.Cart, error) {
	if err := r.db.Save(&cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

func (r *CartRepository) CreateOne(ctx context.Context, cart *entity.Cart) (*entity.Cart, error) {
	if err := r.db.Create(cart).Error; err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *CartRepository) UpdateQuantity(ctx context.Context, id uint, operation string, qty int) error {

	cart := &entity.Cart{}
	if err := r.db.First(&cart, id).Error; err != nil {
		return err
	}
	switch operation {
	case "increase":
		cart.Quantity += qty
	case "decrease":
		if cart.Quantity > 0 {
			cart.Quantity -= qty
		} else {
			return fmt.Errorf("quantity cannot go below 0")
		}
	default:
		return fmt.Errorf("invalid operation, use 'increase' or 'decrease'")
	}

	cart.Subtotal = int(float64(cart.Quantity) * float64(cart.Price))
	if cart.Quantity == 0 {
		if err := r.db.Delete(&cart).Error; err != nil {
			return fmt.Errorf("failed to delete cart: %v", err)
		}
		return nil
	}
	if err := r.db.Save(&cart).Error; err != nil {
		return err
	}
	return nil
}

func (r *CartRepository) DeleteOne(ctx context.Context, id uint) error {
	cart := &entity.Cart{}
	if res := r.db.Model(cart).Where("id = ?", id).Delete(cart); res.Error != nil {
		return res.Error
	}

	return nil
}
