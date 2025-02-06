package entity

type BaseModelOrder struct{}

func (BaseModelOrder) TableName() string {
	return "orders"
}

type Order struct {
	BaseModelOrder
	ID      uint   `json:"id" gorm:"primaryKey"`
	UserID  uint   `json:"user_id"`
	CartID  uint   `json:"cart_id"`
	Total   int    `json:"total" gorm:"not null"`
	Payment bool   `json:"payment" gorm:"default:false"`
	Status  string `json:"status" gorm:"default:pending"`
	User    User   `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Cart    Cart   `json:"cart" gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
}
