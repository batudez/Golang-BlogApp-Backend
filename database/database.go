package database

import (
	"fmt"

	"blogapp.com/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=123456 dbname=BlogDb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Can't connect to the database:", err)
		return
	}
	db.AutoMigrate(&models.Blog{}, &models.User{})
	DB = db
}
