package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
)

type IUserController interface {
	HandleRegister(validator helpers.IValidator, hasher helpers.IHasher) gin.HandlerFunc
}

type UserController struct {
	model models.IUserModel
}

func NewUserController(model models.IUserModel) IUserController {
	return &UserController{
		model: model,
	}
}

func (userController *UserController) HandleRegister(validator helpers.IValidator, hasher helpers.IHasher) gin.HandlerFunc {
	return func(c *gin.Context) {
		// [ ] Memvalidasi request dari json
		// [ ] Memvalidasi apakah email atau attribut unik lain telah terpakai
		// [ ] Melakukan hash pada password
		// [ ] Menyimpan user pada database
		// [ ] Membuat access token
		// [ ] Mengembalikan respon berupa access token
		var user app.User
		user.Photos = []app.Photo{}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		msg, err := validator.Validate(user)
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
		hashedPassword, err := hasher.HashString(user.Password)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Email is already taken",
			})
			return
		}
		user.Password = hashedPassword
		if err := userController.model.CreateUser(&user); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, user)
	}
}
