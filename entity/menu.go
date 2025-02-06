package entity

type BaseModelMenu struct{}

func (BaseModelMenu) TableName() string {
	return "menus"
}

type Menu struct {
	BaseModelMenu
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name" gorm:"unique;not null"`
	Price    int    `json:"price" gorm:"not null"`
	Quantity int    `json:"quantity" gorm:"not null"`
	Image    string `json:"image"`
	User     User   `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
