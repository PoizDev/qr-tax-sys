package controllers

import (
	"fmt"
	"net/http"
	"qrfatura/api/db"
	"qrfatura/api/models"

	"github.com/gin-gonic/gin"
)

// GetAllInvoices returns all invoices with their items and related data
func GetAllInvoices(c *gin.Context) {
	var invoices []models.Fatura
	if err := db.DB.Preload("User").Preload("InvoiceItems.Product").Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Faturalar alınamadı"})
		return
	}
	c.JSON(http.StatusOK, invoices)
}

// GetInvoiceByID returns a single invoice by its ID with full details
func GetInvoiceByID(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Fatura
	if err := db.DB.Preload("User").Preload("InvoiceItems.Product").First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fatura bulunamadı"})
		return
	}
	c.JSON(http.StatusOK, invoice)
}

// GetInvoiceByUserID returns all invoices for a specific user with details
func GetInvoiceByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	var invoices []models.Fatura
	if err := db.DB.Preload("User").Preload("InvoiceItems.Product").Where("user_id = ?", userID).Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı faturaları alınamadı"})
		return
	}
	c.JSON(http.StatusOK, invoices)
}

// AddInvoice creates a new invoice and its items in a transaction, then returns full invoice details
func AddInvoice(c *gin.Context) {
	var req models.Fatura
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek verisi"})
		return
	}

	toplam := 0.0
	for i := range req.InvoiceItems {
		item := &req.InvoiceItems[i]

		var prod models.Product
		if err := db.DB.First(&prod, "product_id = ?", item.ProductID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürün bilgisi alınamadı"})
			return
		}

		item.UnitPrice = fmt.Sprintf("%.2f", prod.UnitPrice)
		item.TaxRate = prod.TaxRate

		subtotal := float64(item.Quantity) * prod.UnitPrice * (1 + prod.TaxRate)
		toplam += subtotal
	}

	req.Total = toplam

	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB transaction başlatılamadı"})
		return
	}

	// Omit nested InvoiceItems to avoid duplicate inserts
	if err := tx.Omit("InvoiceItems").Create(&req).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fatura oluşturulamadı"})
		return
	}

	// Insert each InvoiceItem manually
	for _, item := range req.InvoiceItems {
		item.FaturaID = req.ID
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fatura kalemi oluşturulamadı"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB commit hatası"})
		return
	}

	var fullInv models.Fatura
	if err := db.DB.Preload("User").Preload("InvoiceItems.Product").First(&fullInv, req.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fatura detayları alınamadı"})
		return
	}

	c.JSON(http.StatusOK, fullInv)
}
