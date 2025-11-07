package models

import "time"

type Item struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	SKU        string    `gorm:"unique;not null" json:"sku"`
	CategoryID uint      `json:"category_id"`
	Category   Category  `json:"category,omitempty"`
	Qty        int       `json:"qty"`
	Price      int       `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
