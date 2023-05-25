package models

import (
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
	CreateUser(user *app.UserRegisterRequest) (*User, error)
	GetByEmail(userEmail string, detailed bool) (*User, error)
	GetById(userId uint, detailed bool) (*User, error)
	UpdateUser(user *User, updateBody *app.UserUpdateRequest) (*User, error)
	DeleteUser(user *User) (*User, error)
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

func (userModel *UserModel) GetByEmail(userEmail string, detailed bool) (*User, error) {
	client := userModel.db.GetClient()
	var user User
	var result *gorm.DB
	if detailed {
		result = client.Where("email = ?", userEmail).Preload("Photos").First(&user)
	} else {
		result = client.Where("email = ?", userEmail).First(&user)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (userModel *UserModel) GetById(userId uint, detailed bool) (*User, error) {
	client := userModel.db.GetClient()
	var user User
	var result *gorm.DB
	if detailed {
		result = client.Preload("Photos").First(&user, userId)
	} else {
		result = client.First(&user, userId)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (userModel *UserModel) UpdateUser(user *User, updateBody *app.UserUpdateRequest) (*User, error) {
	client := userModel.db.GetClient()

	user.Username = updateBody.Username

	result := client.Save(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (userModel *UserModel) DeleteUser(user *User) (*User, error) {
	client := userModel.db.GetClient()
	result := client.Delete(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
