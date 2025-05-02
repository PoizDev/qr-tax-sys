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

	// Normal preload kullanımı
	result := db.DB.
		Preload("User").
		Preload("InvoiceItems").
		Preload("InvoiceItems.Product").
		Find(&invoices)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Faturalar alınamadı: " + result.Error.Error()})
		return
	}

	// Veri boş geliyorsa manuel olarak dolduralım
	if len(invoices) > 0 {
		for i := range invoices {
			// Kullanıcı bilgisini manuel doldur
			if invoices[i].UserID > 0 {
				var user models.User
				db.DB.First(&user, invoices[i].UserID)
				invoices[i].User = user
			}

			// Ürün bilgilerini manuel doldur
			for j := range invoices[i].InvoiceItems {
				var product models.Product
				db.DB.First(&product, invoices[i].InvoiceItems[j].ProductID)
				invoices[i].InvoiceItems[j].Product = product
			}
		}
	}

	c.JSON(http.StatusOK, invoices)
}

// GetInvoiceByID returns a single invoice by its ID with full details
func GetInvoiceByID(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Fatura

	// Normal preload kullanımı
	result := db.DB.
		Preload("User").
		Preload("InvoiceItems").
		Preload("InvoiceItems.Product").
		First(&invoice, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fatura bulunamadı: " + result.Error.Error()})
		return
	}

	// Kullanıcı bilgisini manuel doldur
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

// GetInvoiceByUserID returns all invoices for a specific user with details
func GetInvoiceByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	var invoices []models.Fatura

	// Normal preload kullanımı
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

	// Veri boş geliyorsa manuel olarak dolduralım
	if len(invoices) > 0 {
		for i := range invoices {
			// Kullanıcı bilgisini manuel doldur
			if invoices[i].UserID > 0 {
				var user models.User
				db.DB.First(&user, invoices[i].UserID)
				invoices[i].User = user
			}

			// Ürün bilgilerini manuel doldur
			for j := range invoices[i].InvoiceItems {
				var product models.Product
				db.DB.First(&product, invoices[i].InvoiceItems[j].ProductID)
				invoices[i].InvoiceItems[j].Product = product
			}
		}
	}

	c.JSON(http.StatusOK, invoices)
}

// AddInvoice creates a new invoice and its items in a transaction, then returns full invoice details
func AddInvoice(c *gin.Context) {
	var req models.Fatura
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek verisi: " + err.Error()})
		return
	}

	// Kullanıcı kontrolü
	var user models.User
	if err := db.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kullanıcı ID: " + err.Error()})
		return
	}

	// Toplam hesaplama
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

	// Transaction başlatma
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB transaction başlatılamadı: " + tx.Error.Error()})
		return
	}

	// Fatura oluşturma
	if err := tx.Omit("InvoiceItems", "User").Create(&req).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fatura oluşturulamadı: " + err.Error()})
		return
	}

	// Her bir kalemi ekleme
	for i := range req.InvoiceItems {
		item := &req.InvoiceItems[i]
		item.FaturaID = req.ID
		if err := tx.Omit("Product").Create(item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fatura kalemi oluşturulamadı: " + err.Error()})
			return
		}
	}

	// Transaction tamamlama
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB commit hatası: " + err.Error()})
		return
	}

	// Manuel olarak tam veriyi hazırlayalım
	var fullInv models.Fatura
	db.DB.First(&fullInv, req.ID)

	// Kullanıcı bilgisini ekle
	db.DB.First(&fullInv.User, fullInv.UserID)

	// Fatura kalemlerini bul
	var items []models.InvoiceItem
	db.DB.Where("fatura_id = ?", fullInv.ID).Find(&items)
	fullInv.InvoiceItems = items

	// Her bir fatura kalemi için ürün bilgisini ekle
	for i := range fullInv.InvoiceItems {
		var product models.Product
		db.DB.First(&product, fullInv.InvoiceItems[i].ProductID)
		fullInv.InvoiceItems[i].Product = product
	}

	c.JSON(http.StatusOK, fullInv)
}
