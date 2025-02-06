package entity

import "context"

type BaseModelMenu struct{}

func (BaseModelMenu) TableName() string {
	return "menus"
}

type Menu struct {
	BaseModelMenu
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name" gorm:"not null"`
	Price    int    `json:"price" gorm:"not null"`
	Category string `json:"category" gorm:"not null"`
	Quantity int    `json:"quantity" gorm:"not null"`
	Image    string `json:"image"`
	User     User   `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type MenuRepository interface {
	GetMany(ctx context.Context, page, limit int, nameFilter, categoryFilter string) ([]*Menu, int64, error)
	GetOne(ctx context.Context, id uint) (*Menu, error)
	CreateOne(ctx context.Context, menu *Menu) (*Menu, error)
	UpdateOne(ctx context.Context, id uint, updateData map[string]interface{}) (*Menu, error)
	DeleteOne(ctx context.Context, id uint) error
}
