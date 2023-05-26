package models

import (
	"time"

	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
)

type Photo struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Caption   string
	PhotoUrl  string
	UserID    uint
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IPhotoModel interface {
	CreatePhoto(photo *app.PhotoCreationRequest) (*Photo, error)
	GetOwner(userId uint) (*User, error)
	GetAllPhoto() ([]Photo, error)
}

type PhotoModel struct {
	db database.IDatabase
}

func NewPhotoModel(db database.IDatabase) IPhotoModel {
	return &PhotoModel{
		db: db,
	}
}

func (photoModel *PhotoModel) CreatePhoto(photo *app.PhotoCreationRequest) (*Photo, error) {
	newPhoto := Photo{
		Title:    photo.Title,
		Caption:  photo.Caption,
		PhotoUrl: photo.PhotoUrl,
		UserID:   photo.UserID,
	}

	result := photoModel.db.GetClient().Create(&newPhoto)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newPhoto, nil
}

func (photoModel *PhotoModel) GetOwner(userId uint) (*User, error) {
	client := photoModel.db.GetClient()
	var owner User
	result := client.First(&owner, userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &owner, nil
}

func (photoModel *PhotoModel) GetAllPhoto() ([]Photo, error) {
	var photos []Photo
	result := photoModel.db.GetClient().Order("created_at desc").Find(&photos)
	if result.Error != nil {
		return nil, result.Error
	}
	return photos, nil
}
