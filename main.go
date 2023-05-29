package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	app.MaxMultipartMemory = 128 << 20
	db, err := database.CreateMySQLDB(
		os.Getenv("DB_UNAME"),
		os.Getenv("DB_PASSWD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	err = db.MigrateDB(&models.User{}, &models.Photo{})
	if err != nil {
		log.Fatal("Error migrating models to database")
	}
	router.RouteApp(app, db)
	app.Run()
}
