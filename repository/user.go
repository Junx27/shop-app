package repository

import (
	"context"

	"github.com/Junx27/shop-app/entity"
	"gorm.io/gorm"
)

type UserReopository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) entity.UserReopository {
	return &UserReopository{db: db}
}

func (r *UserReopository) GetMany(ctx context.Context, page, limit int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64
	err := r.db.Model(&entity.User{}).Count(&total).Offset((page - 1) * limit).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
