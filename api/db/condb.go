package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("dsn")
	if dsn == "" {
		err := godotenv.Load("db/.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		dsn = os.Getenv("dsn")
	}
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}
