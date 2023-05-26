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

type IPhotoController interface {
	HandleCreatePhoto() gin.HandlerFunc
	Handle
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

// HandleCreatePhoto implements IPhotoController
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
				PhotoUrl: newPhoto.Caption,
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
