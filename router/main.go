package router

import (
	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
)

func RouteApp(app *gin.Engine, database database.IDatabase) {
	UserRouting(app, database)
}
