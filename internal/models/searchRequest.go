package models

import "time"

type SearchRequest struct {
	ID            uint      `gorm:"primaryKey"`
	TransactionID string    `gorm:"index"`
	MessageID     string    `gorm:"index"`
	BapID         string
	BapURI        string
	Domain        string
	City          string
	Intent        string    `gorm:"type:jsonb"`
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
