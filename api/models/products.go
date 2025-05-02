package models

type Product struct {
	ProductID   uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductName string  `json:"product_name" gorm:"not null;"`
	UnitPrice   float64 `json:"unit_price" gorm:"not null;"`
	Photo       []byte  `json:"photo" gorm:"type:mediumblob;default:null;"`
	TaxRate     float64 `json:"tax_rate" gorm:"not null;"`
}
