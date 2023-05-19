package models

import (
	"errors"
	"time"

	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Photos    []Photo
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IUserModel interface {
	CreateUser(user *app.UserRegisterRequest) error
	IsEmailAvailable(userEmail string) bool
}
type UserModel struct {
	db database.IDatabase
}

func NewUserModel(db database.IDatabase) IUserModel {
	return &UserModel{
		db: db,
	}
}

func (userModel *UserModel) CreateUser(u *app.UserRegisterRequest) error {
	newUser := &User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Photos:   []Photo{},
	}

	result := userModel.db.GetClient().Create(newUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (UserModel *UserModel) IsEmailAvailable(userEmail string) bool {
	client := UserModel.db.GetClient()
	result := client.Where("email = ?", userEmail).First(&User{})
	if result.Error != nil {
		return errors.Is(result.Error, gorm.ErrRecordNotFound)
	}
	return false
}
