package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func QRCode(c *gin.Context) {
	id := c.Param("id")
	qrCodeUrl := "192.168.137.27:5000/invoices/" + id

	png, err := qrcode.Encode(qrCodeUrl, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "QR Kod Oluşturulamadı" + err.Error()})
		return
	}

	c.Writer.Header().Set("Content-Type", "image/png")
	c.Writer.Write(png)
}
