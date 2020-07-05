package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/bindings"
	"github.com/ryakosh/wishlist/lib/views"
)

// ErrWishNotFound is returned when Wish does not exist in the database
var ErrWishNotFound = errors.New("Wish not found")

// Wish represents a user's wish to buy something, do something etc.
type Wish struct {
	ID            uint64
	UserID        string
	Name          string `gorm:"type:varchar(256)"`
	Description   string `gorm:"type:varchar(1024)"`
	Link          string
	Image         string
	WantToFulfill []*User `many2many:"want_to_fulfill"`
	Claimers      []*User `many2many:"claimers"`
	Fulfillers    []*User `many2many:"fulfillers"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
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

	db := lib.DB.Create(&wish)
	if db.Error != nil {
		lib.LogError(lib.LPanic, "Could not create wish", db.Error)
	}

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
	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
	} else if db.RecordNotFound() {
		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrWishNotFound,
		}
	}
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

// UpdateWish is used to update wish's general information in the database
func UpdateWish(id uint64, b *bindings.UWish, authedUser string) (*Success, error) {
	var wish Wish

	db := lib.DB.Select("id, user_id").First(&wish, id)
	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
	} else if db.RecordNotFound() {
		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrWishNotFound,
		}
	}

	if wish.UserID != authedUser {
		return nil, &RequestError{
			Status: http.StatusUnauthorized,
			Err:    ErrUserNotAuthorized,
		}
	}

	db = lib.DB.Model(&wish).Updates(&Wish{
		Name:        b.Name,
		Description: b.Description,
		Link:        b.Link,
		Image:       b.Image,
	})
	if db.Error != nil {
		lib.LogError(lib.LPanic, "Could not update wish", db.Error)
	}

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

// DeleteWish is used to delete a wish from the database
func DeleteWish(id uint64, authedUser string) (*Success, error) {
	var wish Wish

	db := lib.DB.Select("id, user_id").First(&wish, id)
	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
	} else if db.RecordNotFound() {
		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrWishNotFound,
		}
	}
	if wish.UserID != authedUser {
		return nil, &RequestError{
			Status: http.StatusUnauthorized,
			Err:    ErrUserNotAuthorized,
		}
	}

	db = lib.DB.Delete(wish)
	if db.Error != nil {
		lib.LogError(lib.LPanic, "Could not delete wish", db.Error)
	}

	return &Success{
		Status: http.StatusOK,
	}, nil
}

func AddWantToFulfill(id uint64, authedUser string) (*Success, error) {
	var wish Wish

	db := lib.DB.Select("id, user_id").First(&wish, id)
	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
	} else if db.RecordNotFound() {
		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrWishNotFound,
		}
	}

	if authedUser == wish.UserID {
		return nil, &RequestError{
			Status: http.StatusUnauthorized,
			Err:    ErrUserNotAuthorized,
		}
	}

	asso := lib.DB.Model(&Wish{ID: id}).Where("user_id = ?", authedUser).Association("WantToFulfill")
	if asso.Error != nil && !gorm.IsRecordNotFoundError(asso.Error) {
		lib.LogError(lib.LPanic, "Could not read wish's WantToFulfill", asso.Error)
	}

	if asso.Count() != 0 {
		return nil, &RequestError{
			Status: http.StatusConflict,
			Err:    ErrUserExists,
		}
	}

	err := lib.DB.Model(&Wish{ID: id}).Association("WantToFulfill").Append(&User{ID: authedUser}).Error
	if err != nil {
		lib.LogError(lib.LPanic, "Could not add to WantToFulfill", err)
	}

	return &Success{
		Status: http.StatusOK,
		View: &views.WishID{
			ID: wish.ID,
		},
	}, nil
}

func init() {
	lib.DB.AutoMigrate(&Wish{})
}
