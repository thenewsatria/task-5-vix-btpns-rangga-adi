package models

import (
	"time"

	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
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
	CreateUser(user *app.UserRegisterRequest) (*User, error)
	GetByEmail(userEmail string) (*User, error)
}
type UserModel struct {
	db database.IDatabase
}

func NewUserModel(db database.IDatabase) IUserModel {
	return &UserModel{
		db: db,
	}
}

func (userModel *UserModel) CreateUser(u *app.UserRegisterRequest) (*User, error) {
	newUser := User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Photos:   []Photo{},
	}

	result := userModel.db.GetClient().Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newUser, nil
}

func (userModel *UserModel) GetByEmail(userEmail string) (*User, error) {
	client := userModel.db.GetClient()
	var user User
	result := client.Where("email = ?", userEmail).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
