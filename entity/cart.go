package entity

import "context"

type BaseModelCart struct{}

func (BaseModelCart) TableName() string {
	return "carts"
}

type Cart struct {
	BaseModelCart
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   uint   `json:"user_id"`
	MenuID   uint   `json:"menu_id"`
	Quantity int    `json:"quantity" gorm:"not null"`
	Subtotal int    `json:"subtotal" gorm:"not null"`
	Status   string `json:"status" gorm:"default:pending"`
	Menu     Menu   `json:"menu" gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE"`
	User     User   `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type CartRepository interface {
	GetOne(ctx context.Context, id uint) (*Cart, error)
	CreateOne(ctx context.Context, cart *Cart) (*Cart, error)
	DeleteOne(ctx context.Context, id uint) error
}
