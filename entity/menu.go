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
	User     User   `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type MenuRepository interface {
	GetMany(ctx context.Context, page, limit int, nameFilter, categoryFilter string) ([]*Menu, int64, error)
	GetOne(ctx context.Context, id uint) (*Menu, error)
	CreateOne(ctx context.Context, menu *Menu) (*Menu, error)
	UpdateOne(ctx context.Context, id uint, updateData map[string]interface{}) (*Menu, error)
	UpdateQuantity(ctx context.Context, id uint, operation string, qty int) error
	DeleteOne(ctx context.Context, id uint) error
}

type MenuService interface {
	CalculateSubTotal(ctx context.Context, id uint, qty int) (int, error)
	DecreaseMenu(ctx context.Context, id uint, qty int) error
	IncreaseMenu(ctx context.Context, id uint, qty int) error
}
