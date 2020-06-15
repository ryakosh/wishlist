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

// CreateWish is used to add a Wish to our database
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

// ReadWish is used to get information about a Wish in our database
func ReadWish(b *bindings.RdWish) (*views.RWish, error) {
	var wish Wish

	db := lib.DB.Omit("fulfilled_by", "created_at", "updated_at").First(&wish, b.ID)

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

// UpdateWish is used to update a Wish in our database
func UpdateWish(b *bindings.UWish, authedUser string) (*views.UWish, error) {
	var wish Wish

	db := lib.DB.Select("id", "user_id").First(&wish, b.ID)
	if !db.RecordNotFound() {
		if wish.UserID == authedUser {
			lib.DB.Model(&wish).Select("name", "description", "link", "image", "fulfilled_by").Updates(&Wish{
				Name:        b.Name,
				Description: b.Description,
				Link:        b.Link,
				Image:       b.Image,
				FulfilledBy: b.FulfilledBy,
			})

			return &views.UWish{
				ID:          b.ID,
				Name:        b.Name,
				Description: b.Description,
				Link:        b.Link,
				Image:       b.Image,
				FulfilledBy: b.FulfilledBy,
			}, nil
		}

		return nil, ErrUserNotAuthorized
	}

	return nil, ErrWishNotFound
}

// DeleteWish is used to delete a Wish from our database
func DeleteWish(b *bindings.RdWish, authedUser string) error {
	var wish Wish

	db := lib.DB.Select("id", "user_id").First(&wish, b.ID)
	if !db.RecordNotFound() {
		if wish.UserID == authedUser {
			lib.DB.Delete(wish)
			return nil
		}

		return ErrUserNotAuthorized
	}

	return ErrWishNotFound
}
