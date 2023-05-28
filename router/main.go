package router

import (
	"github.com/gin-gonic/gin"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
)

func RouteApp(app *gin.Engine, database database.IDatabase) {
	app.Static("/public", "./static/photos")

	UserRouting(app, database)
	PhotoRouting(app, database)
}
