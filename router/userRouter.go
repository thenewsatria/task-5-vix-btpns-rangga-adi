package router

import (
	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/controllers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/middlewares"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
)

func UserRouting(route *gin.Engine, db database.IDatabase) {
	userModel := models.NewUserModel(db)
	validator := helpers.NewValidator()

	userController := controllers.NewUserController(userModel, validator)

	hasher := helpers.NewHasher()
	webToken := helpers.NewWebToken()
	authMW := middlewares.NewAuthMiddleware(userModel, webToken)

	usersRoute := route.Group("/users")
	{
		usersRoute.POST("/register", userController.HandleRegister(hasher, webToken))
		usersRoute.GET("/login", userController.HandleLogin(hasher, webToken))
		idSubRoute := usersRoute.Group("/:userId")
		{
			idSubRoute.Use(authMW.Guard()).Use(authMW.Authorize(userModel))
			{
				idSubRoute.PUT("", userController.HandleUpdate())
				idSubRoute.DELETE("", userController.HandleDelete())
			}
		}
	}
}
