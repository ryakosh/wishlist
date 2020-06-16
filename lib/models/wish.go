package models

import (
	"errors"
	"time"

	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/bindings"
	"github.com/ryakosh/wishlist/lib/views"
)

// ErrWishNotFound is returned when Wish does not exist in the database
var ErrWishNotFound = errors.New("Wish not found")

// Wish represents a user's wish to buy something, do something etc.
type Wish struct {
	ID          uint
	UserID      string
	Name        string `gorm:"type:varchar(256)"`
	Description string `gorm:"type:varchar(1024)"`
	Link        string
	Image       string
	FulfilledBy string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

// CreateWish is used to add a wish to the database
func CreateWish(b *bindings.CWish, authedUser string) *views.CWish {
	wish := Wish{
		UserID:      authedUser,
		Name:        b.Name,
		Description: b.Description,
		Link:        b.Link,
		Image:       b.Image,
	}

	lib.DB.Create(&wish)

	return &views.CWish{
		ID:          wish.ID,
		Name:        wish.Name,
		Description: wish.Description,
		Link:        wish.Link,
		Image:       wish.Image,
	}
}

// ReadWish is used to get general information about a wish in the database
func ReadWish(b *bindings.RdWish) (*views.RWish, error) {
	var wish Wish

	db := lib.DB.Omit("fulfilled_by, created_at, updated_at").First(&wish, b.ID)

	if !db.RecordNotFound() {
		return &views.RWish{
			ID:          wish.ID,
			UserID:      wish.UserID,
			Name:        wish.Name,
			Description: wish.Description,
			Link:        wish.Link,
			Image:       wish.Image,
		}, nil
	}
	return nil, ErrWishNotFound
}

// UpdateWish is used to update wish's general information in the database
func UpdateWish(id uint, b *bindings.UWish, authedUser string) (*views.UWish, error) {
	var wish Wish

	db := lib.DB.Select("id, user_id").First(&wish, id)
	if !db.RecordNotFound() {
		if wish.UserID == authedUser {
			lib.DB.Model(&wish).Updates(&Wish{
				Name:        b.Name,
				Description: b.Description,
				Link:        b.Link,
				Image:       b.Image,
			})

			return &views.UWish{
				ID:          id,
				Name:        b.Name,
				Description: b.Description,
				Link:        b.Link,
				Image:       b.Image,
			}, nil
		}

		return nil, ErrUserNotAuthorized
	}

	return nil, ErrWishNotFound
}

// DeleteWish is used to delete a wish from the database
func DeleteWish(b *bindings.RdWish, authedUser string) error {
	var wish Wish

	db := lib.DB.Select("id, user_id").First(&wish, b.ID)
	if !db.RecordNotFound() {
		if wish.UserID == authedUser {
			lib.DB.Delete(wish)
			return nil
		}

		return ErrUserNotAuthorized
	}

	return ErrWishNotFound
}
