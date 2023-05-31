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
		// [x] Menyimpan user pada database
		// [x] Membuat access token dengan id user yang telah masuk pada database
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

		// Memvalidasi request yang masuk (username, email, password, dsb)
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

		// Mengecek email apakah sudah digunakan oleh user lain atau tidak
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

		// melakukan hashing pada password
		hashedPassword, err := hasher.HashString(registerRequest.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		registerRequest.Password = hashedPassword

		// Membuat user baru pada database sesuai dengan request user
		newUser, err := userController.model.CreateUser(&registerRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// Membuat access token dengan informasi berupa id dari user yang telah dibuat
		accessToken, err := webToken.GenerateToken(newUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// Mengembalikan response berupa json berisi akses token kembali ke client
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
	// [x] Mengambil user terkait dengan email yang diperoleh dari request
	// [x] Melakukan komparasi pada password user saat ini dan password dari request
	// [x] Membuat access token baru dengan informasi berupa id dari user saat ini
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

		// Memvalidasi request json
		msg, _ := userController.validator.Validate(loginRequest)

		if len(msg) != 0 {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}

		// Mengambil user terkait dengan email yang diperoleh dari request
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

		// Melakukan pengecekan password user saat ini (terhash) dengan password dari request (plaintext)
		if !hasher.CheckHash(currentUser.Password, loginRequest.Password) {
			c.JSON(http.StatusUnauthorized, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"message": "Email and password provided doesn't match",
				},
			})
			return
		}

		// membentuk akses token dengan informasi berupa Id dari pengguna saat ini
		accessToken, err := webToken.GenerateToken(currentUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// mengambalikan response berupa akses token kembali ke client
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
		// NOTE: Langkah Kasus Penggunaan User Update
		// [x] Memperoleh user dengan informasi token dari middleware
		// [x] Memvalidasi request json
		// [x] Melakukan pengecekan antara password user saat ini dengan password lama yang dimasukan oleh user
		// [x] Melakukan pengecekan email baru yang dimasukan oleh user
		// [x] Melakukan hashing pada password baru yang dimasukan oleh user
		// [x] Mengupdate user pada database
		// [x] Mengambil informasi tentang photo yang yang terkait dengan user yang diupdate.
		// [x] Membentuk response dari setiap photo yang terkait dengan user yang diupdate.
		// [x] Mengirimkan response kembali ke client.

		// Memperoleh user dari auth middleware
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

		// Memvalidasi request json dari user
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

		// Melakukan pengecekan antara password user saat ini dengan password lama yang dimasukan oleh user
		if !hasher.CheckHash(relatedUser.Password, updateRequest.OldPassword) {
			c.JSON(http.StatusUnauthorized, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"oldPassword": "old password doesn't match the current password.",
				},
			})
			return
		}

		// Melakukan pengecekan email apakah email baru yang dimasukan telah digunakan,
		// Namun apabila email user saat ini sama dengan email yang ada pada request maka proses akan dilanjutkan
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

		// Melakukan hashing pada password baru dari request
		hashedPassword, err := hasher.HashString(updateRequest.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		updateRequest.NewPassword = hashedPassword

		// Melakukan update pada user saat ini dengan informasi sesuai pada request
		updatedUser, err := userController.model.UpdateUser(relatedUser, &updateRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "fail",
				Message: err.Error(),
			})
			return
		}

		// mengambil seluruh photo yang terkait dengan user saat ini.
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

		// membentuk response untuk setiap photo yang ditemukan
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

		// Mengirimkan response kembali ke client
		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.UserDetailGeneralResponse{
				ID:        updatedUser.ID,
				Username:  updatedUser.Username,
				Email:     updatedUser.Email,
				Photos:    &photosResponse,
				CreatedAt: updatedUser.CreatedAt,
				UpdatedAt: updatedUser.UpdatedAt,
			},
		})
	}
}

func (userController *UserController) HandleDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// NOTE: Langkah Kasus Penggunaan User Delete
		// [x] Memperoleh user dengan informasi token dari middleware
		// [x] Mengambil informasi tentang photo yang yang terkait dengan user yang akan dihapus.
		// [x] Menghapus user terkait
		// [x] Menghapus setiap file photo yang terkait dengan user saat ini.
		// [x] Membentuk response dari setiap photo yang terkait dengan user yang dihapus.
		// [x] Mengirimkan response kembali ke client.

		// Memperoleh user dari token pada auth middleware
		relatedUser := c.MustGet("requestedUser").(*models.User)

		// Mengambil user dari database dan juga relasinya dengan photo dengan id dari user
		// yang diperoleh dari middleware
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

		// Menghapus user yang diperoleh dari database
		deletedUser, err := userController.model.DeleteUser(populatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// Menghapus setiap photo yang berkaitan dengan user yang dihapus
		photosResponse := []app.PhotoGeneralResponse{}
		for _, photo := range populatedUser.Photos {

			// menghapus setiap file photo yang berkaitan dengan user yang dihapus
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

			// Membentuk response dari setiap photo yang telah dihapus dari database.
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

		// Mengembalikan response kembali ke client
		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.UserDetailGeneralResponse{
				ID:        deletedUser.ID,
				Username:  deletedUser.Username,
				Email:     deletedUser.Email,
				Photos:    &photosResponse,
				CreatedAt: deletedUser.CreatedAt,
				UpdatedAt: deletedUser.UpdatedAt,
			},
		})
	}
}
