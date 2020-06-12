package models

import (
	"time"
)

// Wish represents a user's wish to buy something, do something etc.
type Wish struct {
	ID          uint       `json:"id"`
	UserID      string     `json:"username"`
	Name        string     `json:"name" gorm:"type:varchar(256)"`
	Description string     `json:"description" gorm:"type:varchar(1024)"`
	Link        string     `json:"link"`
	Image       string     `json:"image"`
	FulfilledBy string     `json:"fulfilled_by"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
