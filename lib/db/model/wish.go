package model

import (
	"errors"
	"time"

	"github.com/ryakosh/wishlist/lib/db"
)

// ErrWishNotFound is returned when Wish does not exist in the database
var ErrWishNotFound = errors.New("Wish not found")

// Wish represents a user's wish to buy something, do something etc.
type Wish struct {
	ID            int
	UserID        string
	Name          string `gorm:"type:varchar(256)"`
	Description   string `gorm:"type:varchar(1024)"`
	Link          string
	Image         string
	WantToFulfill []User `many2many:"want_to_fulfill"`
	Claimers      []User `many2many:"claimers"`
	Fulfillers    []User `many2many:"fulfillers"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

func init() {
	db.DB.AutoMigrate(&Wish{})
}
