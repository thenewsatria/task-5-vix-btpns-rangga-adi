package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/helpers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
	"gorm.io/gorm"
)

type IAuthMiddleware interface {
	Guard() gin.HandlerFunc
	Authorize(model interface{}) gin.HandlerFunc
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
					"token": "There's no token provided, please login",
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
		fmt.Println(tokenStr)
		claims, err := authMW.webToken.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		currentUser, err := authMW.userModel.GetByEmail(claims.Email, false)
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

func (authMW *AuthMiddleware) Authorize(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

		currentUser := c.MustGet("currentUser").(*models.User)
		var ownerId uint = 0
		var castingStat bool = true
		var parseError error = nil
		var queryError error = nil
		var errorResource string = "unknown"

		switch true {
		case c.Param("userId") != "":
			userModel, ok := model.(models.IUserModel)
			if !ok {
				errorResource = "user"
				castingStat = false
				break
			}
			parsedId, err := strconv.ParseUint(c.Param("userId"), 10, 32)
			if err != nil {
				errorResource = "user"
				parseError = err
				break
			}

			requestedUser, err := userModel.GetById(uint(parsedId), false)
			if err != nil {
				errorResource = "user"
				queryError = err
				break
			}
			ownerId = requestedUser.ID
		case c.Param("photoId") != "":
			photoModel, ok := model.(models.IPhotoModel)
			if !ok {
				errorResource = "photo"
				castingStat = false
				break
			}
			parsedId, err := strconv.ParseUint(c.Param("photoId"), 10, 32)
			if err != nil {
				errorResource = "photo"
				parseError = err
				break
			}

			requestedPhoto, err := photoModel.GetById(uint(parsedId), false)
			if err != nil {
				errorResource = "photo"
				queryError = err
				break
			}
			ownerId = requestedPhoto.UserID
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: "There",
			})
			return
		}

		if !castingStat {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: "There's something wrong in the authorization middleware",
			})
			return
		}

		if parseError != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					fmt.Sprintf("%s_id", errorResource): fmt.Sprintf("Invalid %s ID", errorResource),
				},
			})
			return
		}

		if queryError != nil {
			if errors.Is(queryError, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						errorResource: fmt.Sprintf("There's no %s found related with provided %s id", errorResource, errorResource),
					},
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "error",
				Message: queryError.Error(),
			})
			return
		}

		if currentUser.ID == ownerId {
			c.Next()
		} else {
			fmt.Println(queryError)
			c.AbortWithStatusJSON(http.StatusUnauthorized, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					"message": "Access denied, you are unauthorized to access this resource",
				},
			})
			return
		}
	}
}
