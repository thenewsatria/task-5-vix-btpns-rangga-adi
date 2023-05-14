package models

import (
	"errors"

	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
	"gorm.io/gorm"
)

type IUserModel interface {
	CreateUser(user *app.User) error
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

func (userModel *UserModel) CreateUser(user *app.User) error {
	result := userModel.db.GetClient().Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (UserModel *UserModel) IsEmailAvailable(userEmail string) bool {
	client := UserModel.db.GetClient()
	result := client.Where("email = ?", userEmail).First(&app.User{})
	if result.Error != nil {
		return errors.Is(result.Error, gorm.ErrRecordNotFound)
	}
	return false
}
