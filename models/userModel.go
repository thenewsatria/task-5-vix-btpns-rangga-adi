package models

import (
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
)

type IUserModel interface {
	CreateUser(user *app.User) error
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
