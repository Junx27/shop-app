package repository

import (
	"github.com/Junx27/shop-app/entity"
	"gorm.io/gorm"
)

type UserReopository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) entity.UserReopository {
	return &UserReopository{db: db}
}

func (ur *UserReopository) GetMany() ([]*entity.User, error) {
	var users []*entity.User
	if err := ur.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
