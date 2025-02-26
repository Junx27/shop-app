package repository

import (
	"context"

	"github.com/Junx27/shop-app/entity"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func (r *AuthRepository) RegisterUser(ctx context.Context, registerData *entity.User) (*entity.User, error) {
	user := &entity.User{
		UserName:    registerData.UserName,
		Email:       registerData.Email,
		Password:    registerData.Password,
		PhoneNumber: registerData.PhoneNumber,
		Address:     registerData.Address,
	}

	res := r.db.Model(&entity.User{}).Create(user)

	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *AuthRepository) GetUser(ctx context.Context, query interface{}, args ...interface{}) (*entity.User, error) {
	user := &entity.User{}

	if res := r.db.Model(user).Where(query, args...).First(user); res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func NewAuthRepository(db *gorm.DB) entity.AuthRepository {
	return &AuthRepository{
		db: db,
	}
}
