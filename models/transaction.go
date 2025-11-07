package models

import "time"

type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user,omitempty"`
	ItemID    uint      `json:"item_id"`
	Item      Item      `json:"item,omitempty"`
	Qty       int       `json:"qty"`
	Type      string    `json:"type"` // in or out
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}
