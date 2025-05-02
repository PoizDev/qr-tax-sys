package controllers

import (
	"net/http"
	"qrfatura/api/db"
	"qrfatura/api/models"

	"github.com/gin-gonic/gin"
)

//TODO: Fotoğraf eklemek için fonksiyon düzenlenecek.

func GetProductList(c *gin.Context) {
	var products []models.Product
	if err := db.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	productID := c.Param("id")
	var product models.Product
	if err := db.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	var body struct {
		Name      string  `json:"name"`
		UnitPrice float64 `json:"unit_price"`
		TaxRate   float64 `json:"tax_rate"`
		Photo     []byte  `json:"photo"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request " + err.Error()})
		return
	}
	product := models.Product{
		ProductName: body.Name,
		UnitPrice:   body.UnitPrice,
		TaxRate:     body.TaxRate,
		Photo:       body.Photo,
	}

	if err := db.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}
