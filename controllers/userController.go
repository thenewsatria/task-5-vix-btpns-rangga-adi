package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
)

type IUserController interface {
	HandleRegister() gin.HandlerFunc
}

type UserController struct {
	model models.IUserModel
}

func NewUserController(model models.IUserModel) IUserController {
	return &UserController{
		model: model,
	}
}

func (userController *UserController) HandleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user app.User
		user.Photos = []app.Photo{}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := userController.model.CreateUser(&user); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, user)
	}
}
