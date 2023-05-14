package app

import "time"

type Photo struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `json:"userId"`
	User      User `gorm:"constrain:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
