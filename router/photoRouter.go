package router

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/controllers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/middlewares"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
)

func PhotoRouting(route *gin.Engine, db database.IDatabase) {
	photoModel := models.NewPhotoModel(db)
	userModel := models.NewUserModel(db)

	validator := helpers.NewValidator()

	expTime, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		log.Fatal("Error reading token expiration value from .env file")
	}

	webToken := helpers.NewWebToken(expTime, os.Getenv("JWT_SECRET"))

	photoController := controllers.NewPhotoController(photoModel, validator)
	authMW := middlewares.NewAuthMiddleware(userModel, webToken)
	fileUploadMW := middlewares.NewFileUploadMiddleware()
	photoRoute := route.Group("/photos")
	{
		photoRoute.GET("/", photoController.HandleFetchPhotos())
		photoRoute.Use(authMW.Guard())
		{

			photoRoute.POST("", fileUploadMW.AllowMaxSizeKB("photo", 1024), fileUploadMW.AllowedExtension("photo", ".jpeg", ".jpg", ".png"),
				photoController.HandleCreatePhoto())
			idSubRoute := photoRoute.Group("/:photoId")
			{
				idSubRoute.Use(authMW.Authorize(photoModel))
				{
					idSubRoute.PUT("", fileUploadMW.AllowMaxSizeKB("photo", 1024), fileUploadMW.AllowedExtension("photo", ".jpeg", ".jpg", ".png"),
						photoController.HandleUpdatePhoto())
					idSubRoute.DELETE("", photoController.HandleDeletePhoto())
				}
			}
		}
	}
}
