package repository

import (
	"context"
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

func (r *CartRepository) GetMany(ctx context.Context) ([]*entity.Cart, error) {
	var carts []*entity.Cart
	if err := r.db.WithContext(ctx).Find(&carts).Error; err != nil {
		return nil, err
	}
	return carts, nil
}

func (r *CartRepository) GetOne(ctx context.Context, id uint) (*entity.Cart, error) {
	cart := &entity.Cart{}
	if res := r.db.Model(cart).Where("id = ?", id).First(cart); res.Error != nil {
		return nil, res.Error
	}

	return cart, nil
}

// FindByUserAndMenu checks if a cart with the same user and menu already exists.
func (r *CartRepository) FindByUserAndMenu(ctx context.Context, userID uint, menuID uint) (*entity.Cart, error) {
	cart := &entity.Cart{}
	if err := r.db.Where("user_id = ? AND menu_id = ?", userID, menuID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return nil if no cart found
			return nil, nil
		}
		// Return any other error
		return nil, err
	}
	return cart, nil
}

// UpdateOne updates the cart in the database.
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
