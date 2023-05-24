package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
	"gorm.io/gorm"
)

type IUserController interface {
	HandleRegister(hasher helpers.IHasher, webToken helpers.IWebToken) gin.HandlerFunc
	HandleLogin(hasher helpers.IHasher, webToken helpers.IWebToken) gin.HandlerFunc
	HandleUpdate() gin.HandlerFunc
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

		relatedUser, _ := userController.model.GetByEmail(user.Email)
		if relatedUser != nil {
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
		if _, err = userController.model.CreateUser(&user); err != nil {
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

func (userController *UserController) HandleLogin(hasher helpers.IHasher, webToken helpers.IWebToken) gin.HandlerFunc {
	// NOTE: Langkah Kasus Penggunaan User Register
	// [x] Memvalidasi request berupa json
	// [x] Mengambil user terkait dengan email
	// [x] Melakukan komparasi pada password user saat ini dan password dari input pengguna
	// [x] Membuat access token baru
	// [x] Mengembalikan respon berupa access token

	return func(c *gin.Context) {
		var loginRequest app.UserLoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(400, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"json": "Invalid json format",
				},
			})
			return
		}

		msg, err := userController.validator.Validate(loginRequest)
		if err != nil {
			c.JSON(400, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}

		currentUser, err := userController.model.GetByEmail(loginRequest.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(401, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						"message": "Email and password provided doesn't match",
					},
				})
				return
			}
			c.JSON(500, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if !hasher.CheckHash(currentUser.Password, loginRequest.Password) {
			c.JSON(http.StatusUnauthorized, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"message": "Email and password provided doesn't match",
				},
			})
			return
		}

		accessToken, err := webToken.GenerateAccessToken(currentUser.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.UserAuthResponse{
				AccessToken: accessToken,
			},
		})
	}
}

func (UserController *UserController) HandleUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(*models.User)
		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: gin.H{
				"currentEmail": currentUser.Email,
			},
		})
	}
}
