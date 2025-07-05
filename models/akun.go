package models

import (
	"time"
)

type Akun struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string    `gorm:"unique;not null" json:"username"`
	Nohp       string    `json:"nohp"`
	Email      string    `gorm:"unique;not null" json:"email"`
	Password   string    `json:"password"`
	Fotoprofil string    `json:"fotoprofil"`
	UserRole   string    `json:"user_role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Akun) TableName() string {
	return "akun"
}
