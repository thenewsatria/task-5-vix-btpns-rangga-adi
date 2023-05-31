package app

import "time"

type PhotoGeneralResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photoUrl"`
	UserID    uint      `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type FormPhotoCreationRequest struct {
	Title    string `form:"title" valid:"required~title: title is required"`
	Caption  string `form:"caption" valid:"required~caption: caption is required"`
	PhotoUrl string `valid:"required~photoUrl: photo is required please upload a photo"`
	UserID   uint   `valid:"required~userId: userId is required"`
}

type FormPhotoUpdateRequest struct {
	Title    string `form:"title" valid:"required~title: title is required"`
	Caption  string `form:"caption" valid:"required~caption: caption is required"`
	PhotoUrl string `valid:"required~photoUrl: photo is required please upload a photo"`
}

type PhotoDetailGeneralReponse struct {
	ID        uint                 `json:"id"`
	Title     string               `json:"title"`
	Caption   string               `json:"caption"`
	PhotoUrl  string               `json:"photoUrl"`
	Owner     *UserGeneralResponse `json:"owner"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
}
