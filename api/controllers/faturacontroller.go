package controllers

import (
	"fmt"
	"net/http"
	"qrfatura/api/db"
	"qrfatura/api/models"

	"github.com/gin-gonic/gin"
)

func GetAllInvoices(c *gin.Context) {
	var invoices []models.Fatura

	result := db.DB.
		Preload("User").
		Preload("InvoiceItems").
		Preload("InvoiceItems.Product").
		Find(&invoices)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Faturalar alınamadı: " + result.Error.Error()})
		return
	}

	if len(invoices) > 0 {
		for i := range invoices {
			if invoices[i].UserID > 0 {
				var user models.User
				db.DB.First(&user, invoices[i].UserID)
				invoices[i].User = user
			}

			for j := range invoices[i].InvoiceItems {
				var product models.Product
				db.DB.First(&product, invoices[i].InvoiceItems[j].ProductID)
				invoices[i].InvoiceItems[j].Product = product
			}
		}
	}

	c.JSON(http.StatusOK, invoices)
}

func GetInvoiceByID(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Fatura

	result := db.DB.
		Preload("User").
		Preload("InvoiceItems").
		Preload("InvoiceItems.Product").
		First(&invoice, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fatura bulunamadı: " + result.Error.Error()})
		return
	}

	if invoice.UserID > 0 {
		var user models.User
		db.DB.First(&user, invoice.UserID)
		invoice.User = user
	}

	// Ürün bilgilerini manuel doldur
	for i := range invoice.InvoiceItems {
		var product models.Product
		db.DB.First(&product, invoice.InvoiceItems[i].ProductID)
		invoice.InvoiceItems[i].Product = product
	}

	c.JSON(http.StatusOK, invoice)
}

func GetInvoiceByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	var invoices []models.Fatura

	result := db.DB.
		Preload("User").
		Preload("InvoiceItems").
		Preload("InvoiceItems.Product").
		Where("user_id = ?", userID).
		Find(&invoices)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı faturaları alınamadı: " + result.Error.Error()})
		return
	}

	if len(invoices) > 0 {
		for i := range invoices {
			if invoices[i].UserID > 0 {
				var user models.User
				db.DB.First(&user, invoices[i].UserID)
				invoices[i].User = user
			}

			for j := range invoices[i].InvoiceItems {
				var product models.Product
				db.DB.First(&product, invoices[i].InvoiceItems[j].ProductID)
				invoices[i].InvoiceItems[j].Product = product
			}
		}
	}

	c.JSON(http.StatusOK, invoices)
}

func AddInvoice(c *gin.Context) {
	var req models.Fatura
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek verisi: " + err.Error()})
		return
	}

	var user models.User
	if err := db.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID: " + err.Error()})
		return
	}

	toplam := 0.0
	for i := range req.InvoiceItems {
		item := &req.InvoiceItems[i]

		var prod models.Product
		if err := db.DB.First(&prod, item.ProductID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ürün bilgisi alınamadı: " + err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB transaction başlatılamadı: " + tx.Error.Error()})
		return
	}

	if err := tx.Omit("InvoiceItems", "User").Create(&req).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fatura oluşturulamadı: " + err.Error()})
		return
	}

	for i := range req.InvoiceItems {
		item := &req.InvoiceItems[i]
		item.FaturaID = req.ID
		if err := tx.Omit("Product").Create(item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fatura kalemi oluşturulamadı: " + err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB commit hatası: " + err.Error()})
		return
	}

	var fullInv models.Fatura
	db.DB.First(&fullInv, req.ID)

	db.DB.First(&fullInv.User, fullInv.UserID)

	var items []models.InvoiceItem
	db.DB.Where("fatura_id = ?", fullInv.ID).Find(&items)
	fullInv.InvoiceItems = items

	for i := range fullInv.InvoiceItems {
		var product models.Product
		db.DB.First(&product, fullInv.InvoiceItems[i].ProductID)
		fullInv.InvoiceItems[i].Product = product
	}

	c.JSON(http.StatusOK, fullInv)
}

// controllers/faturacontroller.go
func AssignInvoice(c *gin.Context) {
	invID := c.Param("id")
	var body struct {
		UserID uint `json:"user_id"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	var inv models.Fatura
	if err := db.DB.First(&inv, invID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice bulunamadı"})
		return
	}
	inv.UserID = body.UserID
	if err := db.DB.Save(&inv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update hatası"})
		return
	}
	c.JSON(http.StatusOK, inv)
}
