package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
)

type IAuthMiddleware interface {
	Guard() gin.HandlerFunc
}

type AuthMiddleware struct {
	userModel models.IUserModel
	webToken  helpers.IWebToken
}

func NewAuthMiddleware(userModel models.IUserModel, webToken helpers.IWebToken) IAuthMiddleware {
	return &AuthMiddleware{
		userModel: userModel,
		webToken:  webToken,
	}
}

func (authMW *AuthMiddleware) Guard() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"token": "There's no token provided",
				},
			})
			return
		}

		if len(strings.Split(bearerToken, " ")) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"token": "Token provided is invalid",
				},
			})
			return
		}

		tokenStr := strings.Split(bearerToken, " ")[1]
		claims, err := authMW.webToken.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		currentUser, err := authMW.userModel.GetByEmail(claims.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"message": "There's no user found related to the token",
				},
			})
			return
		}
		c.Set("currentUser", currentUser)
		c.Next()
	}
}
