package database

import (
	"fmt"

	"github.com/thenewsatria/task-5-vix-btpns-rangga-adi/app"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IDatabase interface {
	GetClient() *gorm.DB
	Ping() error
	ConnectDB() error
	MigrateDB()
	CloseDB() error
}

type Database struct {
	username string
	password string
	dbname   string
	client   *gorm.DB
}

func CreateMySQLDB(username string, password string, dbname string) (IDatabase, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, dbname)
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db := &Database{username: username, password: password, dbname: dbname, client: client}
	return db, nil
}

func (db *Database) GetClient() *gorm.DB {
	return db.client
}

func (db *Database) Ping() error {
	test, err := db.client.DB()
	if err != nil {
		return err
	}
	err = test.Ping()
	if err != nil {
		return err
	}
	fmt.Println("PINGEDD")
	return nil
}

func (db *Database) ConnectDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db.username, db.password, db.dbname)
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db.client = client
	return nil
}

func (db *Database) MigrateDB() {
	db.client.AutoMigrate(&app.User{}, &app.Photo{})
}

func (db *Database) CloseDB() error {
	test, err := db.client.DB()
	if err != nil {
		return err
	}
	test.Close()
	return nil
}
