package entity

import "context"

type BaseModelOrder struct{}

func (BaseModelOrder) TableName() string {
	return "orders"
}

type OrderReport struct {
	Amount     float64 `json:"amount"`
	PaidAmount float64 `json:"paid_amount"`
	TotalSales float64 `json:"total_sales"`
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
	GetManyByStatus(ctx context.Context, status string) ([]*Order, error)
	UpdatePayment(ctx context.Context, id uint) (*Order, error)
	CreateOne(ctx context.Context, order *Order) (*Order, error)
}

type OrderService interface {
	CalculateOrder(ctx context.Context) (*OrderReport, error)
}
