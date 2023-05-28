package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
	"gorm.io/gorm"
)

type IPhotoController interface {
	HandleCreatePhoto() gin.HandlerFunc
	HandleFetchPhoto() gin.HandlerFunc
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

		currentUser := c.MustGet("currentUser").(*models.User)

		var photoCreationRequest app.FormPhotoCreation
		if err := c.Bind(&photoCreationRequest); err != nil {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"json": "Invalid json format",
				},
			})
			return
		}

		photoCreationRequest.UserID = currentUser.ID
		photoCreationRequest.PhotoUrl = ""

		file, _ := c.FormFile("photo")
		if file != nil {
			timeStamp := time.Now().UnixNano()
			file.Filename = fmt.Sprintf("photos_%d_%d_%s", currentUser.ID, timeStamp, file.Filename)
			photoCreationRequest.PhotoUrl = fmt.Sprintf("http://%s/public/%s", c.Request.Host, file.Filename)
		}

		msg, _ := photoController.validator.Validate(photoCreationRequest)

		if len(msg) != 0 {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}

		newPhoto, err := photoController.model.CreatePhotoForm(&photoCreationRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

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

		err = c.SaveUploadedFile(file, fmt.Sprintf("./static/photos/%s", file.Filename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.PhotoDetailGeneralReponse{
				ID:       newPhoto.ID,
				Title:    newPhoto.Title,
				Caption:  newPhoto.Caption,
				PhotoUrl: newPhoto.PhotoUrl,
				Owner: app.UserGeneralResponse{
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

func (photoController *PhotoController) HandleFetchPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		photos, err := photoController.model.GetAllPhoto()
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
		}

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
		relatedPhoto := c.MustGet("requestedPhoto").(*models.Photo)

		var updateRequest app.PhotoUpdateRequest
		if err := c.ShouldBindJSON(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"json": err.Error(),
				},
			})
			return
		}

		msg, _ := photoController.validator.Validate(updateRequest)

		if len(msg) != 0 {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}

		updatedPhoto, err := photoController.model.UpdatePhoto(relatedPhoto, &updateRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

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

		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.PhotoDetailGeneralReponse{
				ID:       updatedPhoto.ID,
				Title:    updatedPhoto.Title,
				Caption:  updatedPhoto.Caption,
				PhotoUrl: updatedPhoto.PhotoUrl,
				Owner: app.UserGeneralResponse{
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
		relatedPhoto := c.MustGet("requestedPhoto").(*models.Photo)

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

		deletedPhoto, err := photoController.model.DeletePhoto(relatedPhoto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.PhotoDetailGeneralReponse{
				ID:       deletedPhoto.ID,
				Title:    deletedPhoto.Title,
				Caption:  deletedPhoto.Caption,
				PhotoUrl: deletedPhoto.PhotoUrl,
				Owner: app.UserGeneralResponse{
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
