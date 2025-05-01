package models

type InvoiceItem struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	FaturaID  uint     `json:"fatura_id" gorm:"not null;"`
	Fatura    *Fatura  `json:"fatura" gorm:"foreignKey:FaturaID"`
	ProductID uint     `json:"product_id" gorm:"not null;"`
	Product   *Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int      `json:"quantity" gorm:"not null;"`
	UnitPrice string   `json:"unit_price" gorm:"not null;"`
	TaxRate   float64  `json:"tax_rate" gorm:"not null;"`
}
