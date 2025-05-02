package models

import (
	"time"
)

type Fatura struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	UserID       uint          `gorm:"not null" json:"user_id"`
	User         User          `gorm:"foreignKey:UserID" json:"user"` // Pointer olmayan yapı kullanıyoruz
	InvoiceItems []InvoiceItem `gorm:"foreignKey:FaturaID" json:"items"`
	Total        float64       `gorm:"not null" json:"total"`
	CreatedAt    time.Time     `gorm:"autoCreateTime" json:"created_at"`
	Place        string        `gorm:"not null" json:"place"`
}
