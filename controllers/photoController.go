package controllers

import (
	"errors"
	"net/http"
	"strconv"

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

		var photoCreationRequest app.PhotoCreationRequest
		if err := c.ShouldBindJSON(&photoCreationRequest); err != nil {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"json": "Invalid json format",
				},
			})
			return
		}

		photoCreationRequest.UserID = currentUser.ID

		msg, _ := photoController.validator.Validate(photoCreationRequest)

		if len(msg) != 0 {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data:   msg,
			})
			return
		}

		newPhoto, err := photoController.model.CreatePhoto(&photoCreationRequest)
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
						"message": "Can't get owner of the photo, user with related isn't found",
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

		photoId := c.Param("photoId")

		intPhotoId, err := strconv.ParseUint(photoId, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"photo_id": "Invalid photo ID",
				},
			})
			return
		}

		relatedPhoto, err := photoController.model.GetById(uint(intPhotoId), true)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						"photo": "There's no photo found related with provided photo id",
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

		updatedPhoto, err := photoController.model.UpdatePhoto(relatedPhoto, &updateRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
		}

		c.JSON(http.StatusOK, &app.JsendSuccessResponse{
			Status: "success",
			Data: &app.PhotoDetailGeneralReponse{
				ID:       updatedPhoto.ID,
				Title:    updatedPhoto.Title,
				Caption:  updatedPhoto.Caption,
				PhotoUrl: updatedPhoto.PhotoUrl,
				Owner: app.UserGeneralResponse{
					ID:        updatedPhoto.User.ID,
					Username:  updatedPhoto.User.Username,
					Email:     updatedPhoto.User.Email,
					CreatedAt: updatedPhoto.User.CreatedAt,
					UpdatedAt: updatedPhoto.User.UpdatedAt,
				},
				CreatedAt: updatedPhoto.CreatedAt,
				UpdatedAt: updatedPhoto.UpdatedAt,
			},
		})
	}
}
