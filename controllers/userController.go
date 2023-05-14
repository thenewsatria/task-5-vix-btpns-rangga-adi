package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
)

type IUserController interface {
	HandleRegister() gin.HandlerFunc
}

type UserController struct {
	model     models.IUserModel
	validator helpers.IValidator
}

func NewUserController(model models.IUserModel, validator helpers.IValidator) IUserController {
	return &UserController{
		model:     model,
		validator: validator,
	}
}

func (userController *UserController) HandleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user app.User
		user.Photos = []app.Photo{}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		msg, err := userController.validator.Validate(user)
		if err != nil {
			c.JSON(400, gin.H{
				"error": msg,
			})
			return
		}
		if !userController.model.IsEmailAvailable(user.Email) {
			c.JSON(400, gin.H{
				"error": "Email is already taken",
			})
			return
		}
		// TODO: encrypt password using bcrypt here
		if err := userController.model.CreateUser(&user); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, user)
	}
}
