package controllers

import (
	"net/http"
	"qrfatura/api/db"
	"qrfatura/api/models"

	"github.com/gin-gonic/gin"
)

func GetAllItems(c *gin.Context) {
	var items []models.InvoiceItem
	if err := db.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func GetItemsByInvoiceID(c *gin.Context) {
	invoiceID := c.Param("invoice_id")
	var items []models.InvoiceItem
	if err := db.DB.Where("invoice_id = ?", invoiceID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	c.JSON(http.StatusOK, items)
}
