package models

import "time"

type Photo struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User `gorm:"constrain:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
