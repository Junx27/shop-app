package repository

import (
	"context"

	"github.com/Junx27/shop-app/entity"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) entity.CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetOne(ctx context.Context, id uint) (*entity.Cart, error) {
	cart := &entity.Cart{}
	if res := r.db.Model(cart).Where("id = ?", id).First(cart); res.Error != nil {
		return nil, res.Error
	}

	return cart, nil
}

func (r *CartRepository) CreateOne(ctx context.Context, cart *entity.Cart) (*entity.Cart, error) {
	if err := r.db.Create(cart).Error; err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *CartRepository) DeleteOne(ctx context.Context, id uint) error {
	cart := &entity.Cart{}
	if res := r.db.Model(cart).Where("id = ?", id).Delete(cart); res.Error != nil {
		return res.Error
	}

	return nil
}
