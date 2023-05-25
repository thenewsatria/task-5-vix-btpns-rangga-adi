package app

import "time"

type PhotoGeneralResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
