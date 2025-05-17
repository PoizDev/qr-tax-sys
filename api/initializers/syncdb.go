package initializers

import (
	"qrfatura/api/db"
	"qrfatura/api/models"
)

func SyncDB() {
	if db.DB == nil {
		panic("Database connection is not initialized")
	}
	db.DB.AutoMigrate(
		&models.Fatura{},
		&models.User{},
		&models.Product{},
		&models.InvoiceItem{},
	)
}
