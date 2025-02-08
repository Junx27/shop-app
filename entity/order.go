package entity

import "context"

type BaseModelOrder struct{}

func (BaseModelOrder) TableName() string {
	return "orders"
}

type Order struct {
	BaseModelOrder
	ID      uint   `json:"id" gorm:"primaryKey"`
	UserID  uint   `json:"user_id"`
	Total   int    `json:"total" gorm:"not null"`
	Payment bool   `json:"payment" gorm:"default:false"`
	Status  string `json:"status" gorm:"default:pending"`
	User    User   `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type OrderRepository interface {
	GetMany(ctx context.Context, page, limit int) ([]*Order, int64, error)
	CreateOne(ctx context.Context, order *Order) (*Order, error)
}
