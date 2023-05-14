package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/controllers"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/database"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/models"
	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := gin.Default()
	db, err := database.CreateMySQLDB(
		os.Getenv("DB_UNAME"),
		os.Getenv("DB_PASSWD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		panic(err)
	}
	db.MigrateDB()

	userModel := models.NewUserModel(db)

	userController := controllers.NewUserController(userModel)

	router.UserRouting(app, userController)
	app.Run()
}
