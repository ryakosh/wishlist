package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/bindings"
	"github.com/ryakosh/wishlist/lib/views"
)

// ErrWishNotFound is returned when Wish does not exist in the database
var ErrWishNotFound = errors.New("Wish not found")

// Wish represents a user's wish to buy something, do something etc.
type Wish struct {
	ID          uint64
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
func CreateWish(b *bindings.CWish, authedUser string) (*Success, error) {
	wish := Wish{
		UserID:      authedUser,
		Name:        b.Name,
		Description: b.Description,
		Link:        b.Link,
		Image:       b.Image,
	}

	lib.DB.Create(&wish)

	return &Success{
		Status: http.StatusCreated,
		View: &views.CWish{
			ID:          wish.ID,
			Name:        wish.Name,
			Description: wish.Description,
			Link:        wish.Link,
			Image:       wish.Image,
		},
	}, nil
}

// ReadWish is used to get general information about a wish in the database
func ReadWish(id uint64) (*Success, error) {
	var wish Wish

	db := lib.DB.Omit("fulfilled_by, created_at, updated_at").First(&wish, id)

	if !db.RecordNotFound() {
		return &Success{
			Status: http.StatusOK,
			View: &views.RWish{
				ID:          wish.ID,
				UserID:      wish.UserID,
				Name:        wish.Name,
				Description: wish.Description,
				Link:        wish.Link,
				Image:       wish.Image,
			},
		}, nil
	}
	return nil, &RequestError{
		Status: http.StatusNotFound,
		Err:    ErrWishNotFound,
	}
}

// UpdateWish is used to update wish's general information in the database
func UpdateWish(id uint64, b *bindings.UWish, authedUser string) (*Success, error) {
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

			return &Success{
				Status: http.StatusOK,
				View: &views.UWish{
					ID:          id,
					Name:        b.Name,
					Description: b.Description,
					Link:        b.Link,
					Image:       b.Image,
				},
			}, nil
		}

		return nil, &RequestError{
			Status: http.StatusUnauthorized,
			Err:    ErrUserNotAuthorized,
		}
	}

	return nil, &RequestError{
		Status: http.StatusNotFound,
		Err:    ErrWishNotFound,
	}
}

// DeleteWish is used to delete a wish from the database
func DeleteWish(id uint64, authedUser string) (*Success, error) {
	var wish Wish

	db := lib.DB.Select("id, user_id").First(&wish, id)
	if !db.RecordNotFound() {
		if wish.UserID == authedUser {
			lib.DB.Delete(wish)
			return &Success{
				Status: http.StatusOK,
			}, nil
		}

		return nil, &RequestError{
			Status: http.StatusUnauthorized,
			Err:    ErrUserNotAuthorized,
		}
	}

	return nil, &RequestError{
		Status: http.StatusNotFound,
		Err:    ErrWishNotFound,
	}
}

func init() {
	lib.DB.AutoMigrate(&Wish{})
}
