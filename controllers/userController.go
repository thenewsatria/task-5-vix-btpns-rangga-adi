package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
	"gorm.io/gorm"
)

type IUserController interface {
	HandleRegister(hasher helpers.IHasher, webToken helpers.IWebToken) gin.HandlerFunc
	HandleLogin(hasher helpers.IHasher, webToken helpers.IWebToken) gin.HandlerFunc
	HandleUpdate(hasher helpers.IHasher) gin.HandlerFunc
	HandleDelete() gin.HandlerFunc
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

		var registerRequest app.UserRegisterRequest
		if err := c.ShouldBindJSON(&registerRequest); err != nil {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"json": "Invalid json format",
				},
			})
			return
		}

		msg, _ := userController.validator.Validate(registerRequest)

		if registerRequest.Password != registerRequest.ConfirmPassword {
			msg["confirmPassword"] = "password must be matched"
		}

		if len(msg) != 0 {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}

		relatedUser, _ := userController.model.GetByEmail(registerRequest.Email, false)
		if relatedUser != nil {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"email": "Email is already taken",
				},
			})
			return
		}

		hashedPassword, err := hasher.HashString(registerRequest.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		registerRequest.Password = hashedPassword

		newUser, err := userController.model.CreateUser(&registerRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		accessToken, err := webToken.GenerateAccessToken(newUser.ID)
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
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"json": "Invalid json format",
				},
			})
			return
		}

		msg, _ := userController.validator.Validate(loginRequest)

		if len(msg) != 0 {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}

		currentUser, err := userController.model.GetByEmail(loginRequest.Email, false)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						"message": "Email and password provided doesn't match",
					},
				})
				return
			}
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
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

		accessToken, err := webToken.GenerateAccessToken(currentUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.UserAuthResponse{
				AccessToken: accessToken,
			},
		})
	}
}

func (userController *UserController) HandleUpdate(hasher helpers.IHasher) gin.HandlerFunc {
	return func(c *gin.Context) {
		relatedUser := c.MustGet("requestedUser").(*models.User)

		var updateRequest app.UserUpdateRequest
		if err := c.ShouldBindJSON(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"json": "Invalid json format",
				},
			})
			return
		}

		msg, _ := userController.validator.Validate(updateRequest)

		if updateRequest.NewPassword != updateRequest.ConfirmPassword {
			msg["confirmPassword"] = "password must be matched with the new one"
		}

		if len(msg) != 0 {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}

		if !hasher.CheckHash(relatedUser.Password, updateRequest.OldPassword) {
			c.JSON(http.StatusUnauthorized, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"oldPassword": "old password doesn't match the current password.",
				},
			})
			return
		}

		emailOwner, _ := userController.model.GetByEmail(updateRequest.Email, false)
		if emailOwner != nil {
			if emailOwner.Email != relatedUser.Email {
				c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						"email": "Email is already taken",
					},
				})
				return
			}
		}

		hashedPassword, err := hasher.HashString(updateRequest.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		updateRequest.NewPassword = hashedPassword

		updatedUser, err := userController.model.UpdateUser(relatedUser, &updateRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "fail",
				Message: err.Error(),
			})
			return
		}

		populatedUser, err := userController.model.GetById(updatedUser.ID, true)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						"message": "User with related id isn't found",
					},
				})
				return
			}
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		photosResponse := []app.PhotoGeneralResponse{}
		for _, photo := range populatedUser.Photos {
			photosResponse = append(photosResponse, app.PhotoGeneralResponse{
				ID:        photo.ID,
				UserID:    photo.UserID,
				Title:     photo.Title,
				Caption:   photo.Caption,
				PhotoUrl:  photo.PhotoUrl,
				CreatedAt: photo.CreatedAt,
				UpdatedAt: photo.UpdatedAt,
			})
		}

		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.UserDetailGeneralResponse{
				ID:        updatedUser.ID,
				Username:  updatedUser.Username,
				Email:     updatedUser.Email,
				Photos:    photosResponse,
				CreatedAt: updatedUser.CreatedAt,
				UpdatedAt: updatedUser.UpdatedAt,
			},
		})
	}
}

func (userController *UserController) HandleDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		relatedUser := c.MustGet("requestedUser").(*models.User)

		populatedUser, err := userController.model.GetById(relatedUser.ID, true)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						"message": "User with related id isn't found",
					},
				})
				return
			}
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		deletedUser, err := userController.model.DeleteUser(populatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		photosResponse := []app.PhotoGeneralResponse{}
		for _, photo := range populatedUser.Photos {

			// deleting every related photo files.
			strSliceFileLoc := strings.Split(photo.PhotoUrl, "/")
			oldFilename := strSliceFileLoc[len(strSliceFileLoc)-1]
			err = os.Remove(fmt.Sprintf("./static/photos/%s", oldFilename))
			if err != nil {
				c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
					Status:  "error",
					Message: err.Error(),
				})
				return
			}

			photosResponse = append(photosResponse, app.PhotoGeneralResponse{
				ID:        photo.ID,
				UserID:    photo.UserID,
				Title:     photo.Title,
				Caption:   photo.Caption,
				PhotoUrl:  photo.PhotoUrl,
				CreatedAt: photo.CreatedAt,
				UpdatedAt: photo.UpdatedAt,
			})
		}
		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.UserDetailGeneralResponse{
				ID:        deletedUser.ID,
				Username:  deletedUser.Username,
				Email:     deletedUser.Email,
				Photos:    photosResponse,
				CreatedAt: deletedUser.CreatedAt,
				UpdatedAt: deletedUser.UpdatedAt,
			},
		})
	}
}
