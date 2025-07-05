package models

import "time"

type UserFavoritesWisata struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `json:"user_id"`
	WisataId  uint      `json:"wisata_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserFavoritesWisata) TableName() string {
	return "user_favorites_wisata"
}
