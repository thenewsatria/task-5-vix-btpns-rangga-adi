package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
	"gorm.io/gorm"
)

type IPhotoController interface {
	HandleCreatePhoto() gin.HandlerFunc
	HandleFetchPhotos() gin.HandlerFunc
	HandleUpdatePhoto() gin.HandlerFunc
	HandleDeletePhoto() gin.HandlerFunc
}

type PhotoController struct {
	model     models.IPhotoModel
	validator helpers.IValidator
}

func NewPhotoController(model models.IPhotoModel, validator helpers.IValidator) IPhotoController {
	return &PhotoController{
		model:     model,
		validator: validator,
	}
}

func (photoController *PhotoController) HandleCreatePhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		// NOTE: Langkah Kasus Penggunaan Create photo
		// [x] Memperoleh user dengan informasi token dari middleware
		// [x] Mengambil file foto yang diupload dan melakukan rename sehingga unik.
		// [x] Memvalidasi request form-data dari pengguna
		// [x] Membuat data photo pada baru pada database.
		// [x] Mengambil informasi pemilik photo (user) dengan id dari photo yang telah dibuat.
		// [x] Menyimpan file foto yang diupload dan direname ke dalam folder static.
		// [x] Mengirimkan kembali response ke client.

		// Memperoleh user dengan informasi token dari middleware
		currentUser := c.MustGet("currentUser").(*models.User)

		// Melakukan binding antara form-data dari client ke struct
		var photoCreationRequest app.FormPhotoCreationRequest
		if err := c.Bind(&photoCreationRequest); err != nil {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"form-data": "Invalid form-data format",
				},
			})
			return
		}

		// Mengisi nilai UserID dengan id dari user yang saat ini telah terautentikasi.
		photoCreationRequest.UserID = currentUser.ID

		// Inisialisasi nilai PhotoUrl
		photoCreationRequest.PhotoUrl = ""

		// Mengambil informasi file dari form-data dengan key "photo"
		file, _ := c.FormFile("photo")

		// Jika file ada maka rename nama file tersebut sehingga unik.
		if file != nil {
			timeStamp := time.Now().UnixNano()
			file.Filename = fmt.Sprintf("photos_%d_%d_%s", currentUser.ID, timeStamp, file.Filename)
			photoCreationRequest.PhotoUrl = fmt.Sprintf("%s://%s/public/%s", c.Request.URL.Scheme, c.Request.Host, file.Filename)
		}

		// Melakukan validasi pada request
		msg, _ := photoController.validator.Validate(photoCreationRequest)

		if len(msg) != 0 {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}

		// Membuat photo baru pada database sesuai dengan informasi pada request
		newPhoto, err := photoController.model.CreatePhoto(&photoCreationRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// Mengambil informasi mengenai pemilik photo (user) dengan id dari photo yang baru saja dibuat.
		photoOwner, err := photoController.model.GetOwner(newPhoto.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						"message": "Can't populate owner of the photo, user with related userId is not found",
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

		// Menyimpan file yang diupload pada folder static
		err = c.SaveUploadedFile(file, fmt.Sprintf("./static/photos/%s", file.Filename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// Mengirimkan kembali response ke client.
		c.JSON(http.StatusCreated, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.PhotoDetailGeneralReponse{
				ID:       newPhoto.ID,
				Title:    newPhoto.Title,
				Caption:  newPhoto.Caption,
				PhotoUrl: newPhoto.PhotoUrl,
				Owner: &app.UserGeneralResponse{
					ID:        photoOwner.ID,
					Username:  photoOwner.Username,
					Email:     photoOwner.Email,
					CreatedAt: photoOwner.CreatedAt,
					UpdatedAt: photoOwner.UpdatedAt,
				},
				CreatedAt: newPhoto.CreatedAt,
				UpdatedAt: newPhoto.UpdatedAt,
			},
		})
	}
}

func (photoController *PhotoController) HandleFetchPhotos() gin.HandlerFunc {
	return func(c *gin.Context) {
		// NOTE: Langkah Kasus Penggunaan Fetch photos
		// [x] Mengambil seluruh data photo dari database
		// [x] Membentuk response untuk masing masing photo
		// [x] Mengirimkan response kembali ke client.

		// Mengambil semua photo yang terdapat pada database
		photos, err := photoController.model.GetAllPhoto()
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// Membentuk response untuk masing masing photo yang diperoleh
		photosReponse := []*app.PhotoGeneralResponse{}
		for _, photo := range photos {
			photosReponse = append(photosReponse, &app.PhotoGeneralResponse{
				ID:        photo.ID,
				Title:     photo.Title,
				Caption:   photo.Caption,
				PhotoUrl:  photo.PhotoUrl,
				UserID:    photo.UserID,
				CreatedAt: photo.CreatedAt,
				UpdatedAt: photo.UpdatedAt,
			})
		}

		// Mengirimkan response kembali ke client
		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: gin.H{
				"photos": photosReponse,
			},
		})
	}
}

// HandleUpdatePhoto implements IPhotoController
func (photoController *PhotoController) HandleUpdatePhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		// NOTE: Langkah Kasus Penggunaan Update photo
		// [x] Memperoleh photo dengan photo id dari middleware authorization
		// [x] Memperoleh user dengan informasi token dari middleware
		// [x] Memperoleh nama file dengan photoUrl dari photo yang akan diupdate.
		// [x] Mengambil file foto yang diupload dan melakukan rename sehingga unik.
		// [x] Melakukan validasi pada data yang diberikan pengguna
		// [x] Melakukan update pada photo dengan data dari request
		// [x] Memperoleh informasi dari pemilik photo (user) dengan id dari photo yang baru saja diupdate
		// [x] Simpan file yang diupload (jika ada)
		// [x] Hapus file lama jika ada file baru yang diupload
		// [x] Mengirimkan kembali response ke client.

		// Memperoleh photo dengan photo id dari middleware authorization
		relatedPhoto := c.MustGet("requestedPhoto").(*models.Photo)

		// Memperoleh user dengan informasi token dari middleware
		currentUser := c.MustGet("currentUser").(*models.User)

		var photoUpdateRequest app.FormPhotoUpdateRequest
		if err := c.Bind(&photoUpdateRequest); err != nil {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"form-data": "Invalid form-data format",
				},
			})
			return
		}

		// Memperoleh nama file dengan photoUrl dari photo yang akan diupdate
		photoUpdateRequest.PhotoUrl = relatedPhoto.PhotoUrl
		strSliceFileLoc := strings.Split(relatedPhoto.PhotoUrl, "/")
		oldFilename := strSliceFileLoc[len(strSliceFileLoc)-1]

		// Memperoleh informasi dari file yang diupload
		file, _ := c.FormFile("photo")

		// Melakukan perubahan nama file yang diupload sehingga bersifat unik
		if file != nil {
			timeStamp := time.Now().UnixNano()
			file.Filename = fmt.Sprintf("photos_%d_%d_%s", currentUser.ID, timeStamp, file.Filename)
			photoUpdateRequest.PhotoUrl = fmt.Sprintf("%s://%s/public/%s", c.Request.URL.Scheme, c.Request.Host, file.Filename)
		}

		// Melakukan validasi pada data yang diberikan pengguna
		msg, _ := photoController.validator.Validate(photoUpdateRequest)

		if len(msg) != 0 {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}

		// Melakukan update pada photo dengan data dari request
		updatedPhoto, err := photoController.model.UpdatePhoto(relatedPhoto, &photoUpdateRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// Memperoleh informasi dari pemilik photo (user) dengan id dari photo yang baru saja diupdate
		photoOwner, err := photoController.model.GetOwner(updatedPhoto.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						"message": "Can't populate owner of the photo, user with related userId is not found",
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

		// jika terdapat file yang diupload maka simpan file tersebut pada folder static dan hapus file lama
		if file != nil {

			// Menyimpan file baru.
			err = c.SaveUploadedFile(file, fmt.Sprintf("./static/photos/%s", file.Filename))
			if err != nil {
				c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
					Status:  "error",
					Message: err.Error(),
				})
				return
			}

			// Menghapus file lama
			err = os.Remove(fmt.Sprintf("./static/photos/%s", oldFilename))
			if err != nil {
				c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
					Status:  "error",
					Message: err.Error(),
				})
				return
			}
		}

		// Mengirimkan response kembali ke client.
		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.PhotoDetailGeneralReponse{
				ID:       updatedPhoto.ID,
				Title:    updatedPhoto.Title,
				Caption:  updatedPhoto.Caption,
				PhotoUrl: updatedPhoto.PhotoUrl,
				Owner: &app.UserGeneralResponse{
					ID:        photoOwner.ID,
					Username:  photoOwner.Username,
					Email:     photoOwner.Email,
					CreatedAt: photoOwner.CreatedAt,
					UpdatedAt: photoOwner.UpdatedAt,
				},
				CreatedAt: updatedPhoto.CreatedAt,
				UpdatedAt: updatedPhoto.UpdatedAt,
			},
		})
	}
}

func (photoController *PhotoController) HandleDeletePhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		// NOTE: Langkah Kasus Penggunaan Delete photo
		// [x] Memperoleh photo dengan photo id dari middleware authorization
		// [x] Mendapatkan informasi mengenai pemilik dari photo
		// [x] Menghapus data photo dari database
		// [x] Menghapus file photo yang terkait dengan photo yang dihapus
		// [x] Mengirim response kembali ke client.

		// Mengambil photo dengan photo id dari middleware authorization
		relatedPhoto := c.MustGet("requestedPhoto").(*models.Photo)

		// Mengambil pemilik dari photo (user) dengan photo id
		photoOwner, err := photoController.model.GetOwner(relatedPhoto.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						"message": "Can't populate owner of the photo, user with related userId is not found",
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

		// Menghapus photo terkait dari database
		deletedPhoto, err := photoController.model.DeletePhoto(relatedPhoto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// // Menghapus file photo yang terkait dengan photo yang dihapus
		strSliceFileLoc := strings.Split(deletedPhoto.PhotoUrl, "/")
		oldFilename := strSliceFileLoc[len(strSliceFileLoc)-1]
		err = os.Remove(fmt.Sprintf("./static/photos/%s", oldFilename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// Mengirim response kembali ke client
		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.PhotoDetailGeneralReponse{
				ID:       deletedPhoto.ID,
				Title:    deletedPhoto.Title,
				Caption:  deletedPhoto.Caption,
				PhotoUrl: deletedPhoto.PhotoUrl,
				Owner: &app.UserGeneralResponse{
					ID:        photoOwner.ID,
					Username:  photoOwner.Username,
					Email:     photoOwner.Email,
					CreatedAt: photoOwner.CreatedAt,
					UpdatedAt: photoOwner.UpdatedAt,
				},
				CreatedAt: deletedPhoto.CreatedAt,
				UpdatedAt: deletedPhoto.UpdatedAt,
			},
		})
	}
}
