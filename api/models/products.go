package models

type Product struct {
	ProductID   uint    `json:"id" gorm:"primaryKey"`
	ProductName string  `json:"product_name" gorm:"not null;"`
	UnitPrice   string  `json:"unit_price" gorm:"not null;"`
	Photo       []byte  `json:"photo" gorm:"not null;type:mediumblob;"`
	TaxRate     float64 `json:"tax_rate" gorm:"not null;"`
}
