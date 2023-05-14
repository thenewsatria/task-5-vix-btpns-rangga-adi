package router

import (
	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/controllers"
)

func UserRouting(route *gin.Engine, userController controllers.IUserController) {
	usersRoute := route.Group("/users")
	{
		usersRoute.POST("/register", userController.HandleRegister())
	}
}
