package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
)

type IUserController interface {
	HandleRegister(hasher helpers.IHasher, webToken helpers.IWebToken) gin.HandlerFunc
}

type UserController struct {
	model     models.IUserModel
	validator helpers.IValidator
}

func NewUserController(model models.IUserModel, validator helpers.IValidator) IUserController {
	return &UserController{
		model:     model,
		validator: validator,
	}
}

func (userController *UserController) HandleRegister(hasher helpers.IHasher, webToken helpers.IWebToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		// NOTE: Langkah Kasus Penggunaan User Register
		// [x] Memvalidasi request berupa json
		// [x] Memvalidasi apakah email atau attribut unik lain telah terpakai
		// [x] Melakukan hash pada password
		// [x] Membuat access token
		// [x] Menyimpan user pada database
		// [x] Mengembalikan respon berupa access token

		var user app.UserRegisterRequest
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"json": "Invalid json format",
				},
			})
			return
		}
		msg, err := userController.validator.Validate(user)
		if err != nil {
			c.JSON(400, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}
		if !userController.model.IsEmailAvailable(user.Email) {
			c.JSON(400, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"email": "Email is already taken",
				},
			})
			return
		}
		hashedPassword, err := hasher.HashString(user.Password)
		if err != nil {
			c.JSON(500, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		user.Password = hashedPassword
		accessToken, err := webToken.GenerateAccessToken(user.Email)
		if err != nil {
			c.JSON(500, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		if err := userController.model.CreateUser(&user); err != nil {
			c.JSON(500, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		c.JSON(201, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.UserAuthResponse{
				AccessToken: accessToken,
			},
		})
	}
}
