package router

import (
	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/controllers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
)

func UserRouting(route *gin.Engine, userController controllers.IUserController) {
	validator := &helpers.Validator{}
	hasher := &helpers.Hasher{}

	usersRoute := route.Group("/users")
	{
		usersRoute.POST("/register", userController.HandleRegister(validator, hasher))
	}
}
