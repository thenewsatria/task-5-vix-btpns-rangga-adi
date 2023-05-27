package router

import (
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

	webToken := helpers.NewWebToken()

	photoController := controllers.NewPhotoController(photoModel, validator)
	authMW := middlewares.NewAuthMiddleware(userModel, webToken)

	photoRoute := route.Group("/photos")
	{
		photoRoute.GET("/", photoController.HandleFetchPhoto())
		photoRoute.Use(authMW.Guard())
		{
			photoRoute.POST("/", photoController.HandleCreatePhoto())
			idSubRoute := photoRoute.Group("/:photoId")
			{
				idSubRoute.Use(authMW.Authorize(photoModel))
				{
					idSubRoute.PUT("/", photoController.HandleUpdatePhoto())
				}
			}
		}
	}
}
