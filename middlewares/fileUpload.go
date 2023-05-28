package middlewares

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
)

type IFileUploadMiddleware interface {
	AllowedExtension(formName string, fileExtensions ...string) gin.HandlerFunc
	AllowMaxSizeKB(formName string, maxSizeKB uint64) gin.HandlerFunc
}

type FileUploadMiddleware struct {
}

func NewFileUploadMiddleware() IFileUploadMiddleware {
	return &FileUploadMiddleware{}
}

// AllowedExt implements IFileUploadMiddleware
func (fileUploadMW *FileUploadMiddleware) AllowedExtension(formName string, fileExtensions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile(formName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &app.JsendErrorResponse{
				Status:  "fail",
				Message: err.Error(),
			})
			return
		}
		if file != nil {
			fmt.Println(file.Filename)
		}
		if file == nil {
			c.Next()
			return
		} else {
			for _, ext := range fileExtensions {
				if strings.ToLower(ext) == filepath.Ext(file.Filename) {
					c.Next()
					return
				}
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, &app.JsendFailResponse{
				Status: "fail",
				Data: gin.H{
					formName: fmt.Sprintf("%s with %s extension is not allowed", formName, filepath.Ext(file.Filename)),
				},
			})
			return
		}
	}
}

// AllowMaxSizeKb implements IFileUploadMiddleware
func (fileUploadMW *FileUploadMiddleware) AllowMaxSizeKB(formName string, maxSizeKB uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("masuk checksize")
		file, _ := c.FormFile(formName)
		if file != nil {
			fmt.Println(file.Filename)
		}
		if file == nil {
			fmt.Println("file is nil")
			c.Next()
			return
		} else {
			if file.Size <= int64(maxSizeKB)*1024 {
				fmt.Println("masuk bos 123123")
				c.Next()
				return
			} else {
				fmt.Println("masuk bos")
				c.AbortWithStatusJSON(http.StatusBadRequest, &app.JsendFailResponse{
					Status: "fail",
					Data: gin.H{
						formName: fmt.Sprintf("%s file is to large, the %s is larger than %dKB.", formName, formName, maxSizeKB),
					},
				})
				return
			}
		}
	}
}
