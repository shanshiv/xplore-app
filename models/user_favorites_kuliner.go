package models

import "time"

type UserFavoritesKuliner struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `json:"user_id"`
	KulinerId uint      `json:"kuliner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserFavoritesKuliner) TableName() string {
	return "user_favorites_kuliner"
}
