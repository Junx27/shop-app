package entity

import (
	"context"

	"gorm.io/gorm"
)

type BaseModelCart struct{}

func (BaseModelCart) TableName() string {
	return "carts"
}

type Cart struct {
	BaseModelCart
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   uint   `json:"user_id"`
	MenuID   uint   `json:"menu_id"`
	OrderID  *uint  `json:"order_id"`
	Quantity int    `json:"quantity" gorm:"not null"`
	Subtotal int    `json:"subtotal" gorm:"not null"`
	Status   string `json:"status" gorm:"default:pending"`
	Menu     Menu   `json:"-" gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE"`
	User     User   `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Order    Order  `json:"-" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

func (cart *Cart) BeforeSave(tx *gorm.DB) (err error) {
	if cart.Status != "checkout" {
		cart.OrderID = nil
	}
	return nil
}

type CartRepository interface {
	GetMany(ctx context.Context) ([]*Cart, error)
	GetOne(ctx context.Context, id uint) (*Cart, error)
	CreateOne(ctx context.Context, cart *Cart) (*Cart, error)
	UpdateQuantity(ctx context.Context, id uint, operation string, qty int) error
	DeleteOne(ctx context.Context, id uint) error
	FindByUserAndMenu(ctx context.Context, userID uint, menuID uint) (*Cart, error)
	UpdateOne(ctx context.Context, cart *Cart) (*Cart, error)
}

type CartService interface {
	DecreaseCart(ctx context.Context, id uint, qty int) error
	IncreaseCart(ctx context.Context, id uint, qty int) error
}
