package repository

import (
	"context"
	"fmt"

	"github.com/Junx27/shop-app/entity"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) entity.MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) GetMany(ctx context.Context, page, limit int, nameFilter, categoryFilter string) ([]*entity.Menu, int64, error) {
	var menus []*entity.Menu
	var totalItems int64
	query := r.db.Model(&entity.Menu{})

	if nameFilter != "" {
		query = query.Where("name LIKE ?", "%"+nameFilter+"%")
	}

	if categoryFilter != "" {
		query = query.Where("category LIKE ?", "%"+categoryFilter+"%")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Preload("User").Find(&menus).Error; err != nil {
		return nil, 0, err
	}

	return menus, totalItems, nil
}

func (r *MenuRepository) GetOne(ctx context.Context, id uint) (*entity.Menu, error) {
	menu := &entity.Menu{}
	if res := r.db.Model(menu).Where("id = ?", id).Preload("User").First(menu); res.Error != nil {
		return nil, res.Error
	}

	return menu, nil
}

func (r *MenuRepository) CreateOne(ctx context.Context, menu *entity.Menu) (*entity.Menu, error) {
	if err := r.db.WithContext(ctx).Create(menu).Error; err != nil {
		return nil, err
	}

	return menu, nil
}

func (r *MenuRepository) UpdateOne(ctx context.Context, id uint, updateData map[string]interface{}) (*entity.Menu, error) {
	menu := &entity.Menu{}
	res := r.db.Model(&menu).Where("id = ?", id).Updates(updateData)

	if res.Error != nil {
		return nil, res.Error
	}

	return menu, nil
}
func (r *MenuRepository) UpdateQuantity(ctx context.Context, id uint, operation string, qty int) error {

	menu := &entity.Menu{}
	if err := r.db.First(&menu, id).Error; err != nil {
		return err
	}
	switch operation {
	case "increase":
		menu.Quantity += qty
	case "decrease":
		if menu.Quantity > 0 {
			menu.Quantity -= qty
		} else {
			return fmt.Errorf("quantity cannot go below 0")
		}
	default:
		return fmt.Errorf("invalid operation, use 'increase' or 'decrease'")
	}
	if err := r.db.Save(&menu).Error; err != nil {
		return err
	}
	return nil
}

func (r *MenuRepository) DeleteOne(ctx context.Context, id uint) error {
	menu := &entity.Menu{}
	res := r.db.Model(&menu).Where("id = ?", id).Delete(menu)

	if res.Error != nil {
		return res.Error
	}

	return nil
}
