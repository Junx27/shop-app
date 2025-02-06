package entity

import "gorm.io/gorm"

type UserRole string

const (
	Admin UserRole = "admin"
)

type BaseModelUser struct{}

func (BaseModelUser) TableName() string {
	return "users"
}

type User struct {
	BaseModelUser
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserName    string `json:"username" gorm:"unique;not null"`
	Email       string `json:"email" gorm:"unique;not null"`
	Password    string `json:"password" gorm:"not null"`
	PhoneNumber int    `json:"phone_number"`
	Role        string `json:"role" gorm:"default:user"`
}

type UserResponse struct {
	BaseModelUser
	ID       uint   `json:"-"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type UserReopository interface {
	GetMany() ([]*User, error)
}

func (u *User) AfterCreate(db *gorm.DB) (err error) {
	if u.ID == 1 {
		db.Model(u).Update("role", Admin)
	}
	return
}
