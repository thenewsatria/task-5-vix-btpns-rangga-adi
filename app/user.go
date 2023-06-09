package app

import "time"

type UserRegisterRequest struct {
	Username        string `json:"username" valid:"required~username: username is required"`
	Email           string `json:"email" valid:"email,required~email: email is required"`
	Password        string `json:"password" valid:"required~password: password is required,minstringlength(6)~password: password must be at least 6 characters"`
	ConfirmPassword string `json:"confirmPassword" valid:"required~confirmPassword: confirm password is required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" valid:"email,required~email: email is required"`
	Password string `json:"password" valid:"required~password: password is required"`
}

type UserUpdateRequest struct {
	Username        string `json:"username" valid:"required~username: username is required"`
	Email           string `json:"email" valid:"email,required~email: email is required"`
	OldPassword     string `json:"oldPassword" valid:"required~oldPassword: old password password is required"`
	NewPassword     string `json:"newPassword" valid:"required~newPassword: new password is required,minstringlength(6)~newPassword: new password must be at least 6 characters"`
	ConfirmPassword string `json:"confirmPassword" valid:"required~confirmPassword: confirm password is required"`
}

type UserDetailGeneralResponse struct {
	ID        uint                    `json:"id"`
	Username  string                  `json:"username"`
	Email     string                  `json:"email"`
	Photos    *[]PhotoGeneralResponse `json:"photos"`
	CreatedAt time.Time               `json:"createdAt"`
	UpdatedAt time.Time               `json:"updatedAt"`
}

type UserGeneralResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserAuthResponse struct {
	AccessToken string `json:"accessToken"`
}
