package app

// type User struct {
// 	ID        uint      `gorm:"primaryKey" json:"id"`
// 	Username  string    `gorm:"not null" json:"username" valid:"required~username: username is required"`
// 	Email     string    `gorm:"unique;not null" json:"email" valid:"email,required~email: email is required"`
// 	Password  string    `gorm:"not null" json:"password" valid:"required~password: password is required,minstringlength(6)~password: password must be at least 6 characters"`
// 	Photos    []Photo   `json:"photos"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }

type UserRegisterRequest struct {
	Username string `json:"username" valid:"required~username: username is required"`
	Email    string `json:"email" valid:"email,required~email: email is required"`
	Password string `json:"password" valid:"required~password: password is required,minstringlength(6)~password: password must be at least 6 characters"`
}

type UserLoginRequest struct {
	Email    string `json:"email" valid:"email,required~email: email is required"`
	Password string `json:"password" valid:"required~password: password is required"`
}

type UserUpdateRequest struct {
	Username string `json:"username" valid:"required~username: username is required"`
	Password string `json:"password" valid:"required~password: password is required,minstringlength(6)~password: password must be at least 6 characters"`
}

type UserAuthResponse struct {
	AccessToken string `json:"accessToken"`
}
