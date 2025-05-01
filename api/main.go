package main

import (
	"qrfatura/api/controllers"
	"qrfatura/api/db"
	"qrfatura/api/initializers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	// Initialize the database connection
	db.Connect()

	// Sync the database
	initializers.SyncDB()
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/users", controllers.GetUserList)

	if err := r.Run(":5000"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
