package models

import (
	"time"

	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
	"gorm.io/gorm"
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
	CreatePhotoForm(photo *app.FormPhotoCreation) (*Photo, error)
	GetAllPhoto() ([]Photo, error)
	GetOwner(userId uint) (*User, error)
	GetById(photoId uint, detailed bool) (*Photo, error)
	UpdatePhoto(photo *Photo, updateBody *app.PhotoUpdateRequest) (*Photo, error)
	UpdatePhotoForm(photo *Photo, updateBody *app.FormPhotoUpdateRequest) (*Photo, error)
	DeletePhoto(photo *Photo) (*Photo, error)
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

func (photoModel *PhotoModel) CreatePhotoForm(photo *app.FormPhotoCreation) (*Photo, error) {
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

func (photoModel *PhotoModel) UpdatePhoto(photo *Photo, updateBody *app.PhotoUpdateRequest) (*Photo, error) {
	client := photoModel.db.GetClient()

	photo.Title = updateBody.Title
	photo.Caption = updateBody.Caption
	photo.PhotoUrl = updateBody.PhotoUrl

	result := client.Save(photo)
	if result.Error != nil {
		return nil, result.Error
	}

	return photo, nil
}

func (photoModel *PhotoModel) UpdatePhotoForm(photo *Photo, updateBody *app.FormPhotoUpdateRequest) (*Photo, error) {
	client := photoModel.db.GetClient()

	photo.Title = updateBody.Title
	photo.Caption = updateBody.Caption
	photo.PhotoUrl = updateBody.PhotoUrl

	result := client.Save(photo)
	if result.Error != nil {
		return nil, result.Error
	}

	return photo, nil
}

func (photoModel *PhotoModel) GetById(photoId uint, detailed bool) (*Photo, error) {
	client := photoModel.db.GetClient()

	var photo Photo
	var result *gorm.DB
	if detailed {
		result = client.Preload("User").First(&photo, photoId)
	} else {
		result = client.First(&photo, photoId)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &photo, nil
}

func (photoModel *PhotoModel) DeletePhoto(photo *Photo) (*Photo, error) {
	client := photoModel.db.GetClient()
	result := client.Delete(photo)
	if result.Error != nil {
		return nil, result.Error
	}

	return photo, nil
}
