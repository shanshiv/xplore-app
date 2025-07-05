package models

import "time"

type UserFavoritesPenginapan struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `json:"user_id"`
	PenginapanId uint      `json:"penginapan_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (UserFavoritesPenginapan) TableName() string {
	return "user_favorites_penginapan"
}
