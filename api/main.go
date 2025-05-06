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
	//'User Controller Routes
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/users", controllers.GetUserList)

	//' Prod Controller Routes
	r.GET("/products", controllers.GetProductList)
	r.GET("/products/:id", controllers.GetProductByID)
	r.POST("/products/create", controllers.CreateProduct)

	//' Item Controller Routes
	r.GET("/items", controllers.GetAllItems)
	r.GET("/items/:invoice_id", controllers.GetItemsByInvoiceID)
	//' Invoice Controller Routes
	r.GET("/invoices", controllers.GetAllInvoices)
	r.GET("/invoices/:id", controllers.GetInvoiceByID)
	r.GET("/invoices/user/:user_id", controllers.GetInvoiceByUserID)
	r.POST("/invoices/create", controllers.AddInvoice)

	//' QR Code Controller Routes
	r.GET("/qrcode/:id", controllers.QRCode)

	if err := r.Run(":5000"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
