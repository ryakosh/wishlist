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
	ID            uint64
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

// CreateWish is used to add a wish to the database
// func CreateWish(o *Options) (*Success, error) {
// 	b := o.B.(*bindings.CWish)

// 	wish := Wish{
// 		UserID:      o.AuthedUser,
// 		Name:        b.Name,
// 		Description: b.Description,
// 		Link:        b.Link,
// 		Image:       b.Image,
// 	}

// 	db := lib.DB.Create(&wish)
// 	if db.Error != nil {
// 		lib.LogError(lib.LPanic, "Could not create wish", db.Error)
// 	}

// 	return &Success{
// 		Status: http.StatusCreated,
// 		View: &views.CWish{
// 			ID:          wish.ID,
// 			Name:        wish.Name,
// 			Description: wish.Description,
// 			Link:        wish.Link,
// 			Image:       wish.Image,
// 		},
// 	}, nil
// }

// // ReadWish is used to get general information about a wish in the database
// func ReadWish(o *Options) (*Success, error) {
// 	var wish Wish

// 	id := o.Params["id"].(uint64)

// 	db := lib.DB.Omit("fulfilled_by, created_at, updated_at").First(&wish, id)
// 	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
// 		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
// 	} else if db.RecordNotFound() {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrWishNotFound,
// 		}
// 	}
// 	return &Success{
// 		Status: http.StatusOK,
// 		View: &views.RWish{
// 			ID:          wish.ID,
// 			UserID:      wish.UserID,
// 			Name:        wish.Name,
// 			Description: wish.Description,
// 			Link:        wish.Link,
// 			Image:       wish.Image,
// 		},
// 	}, nil
// }

// // UpdateWish is used to update wish's general information in the database
// func UpdateWish(o *Options) (*Success, error) {
// 	var wish Wish

// 	id := o.Params["id"].(uint64)
// 	b := o.B.(*bindings.UWish)

// 	db := lib.DB.Select("id, user_id").First(&wish, id)
// 	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
// 		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
// 	} else if db.RecordNotFound() {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrWishNotFound,
// 		}
// 	}

// 	if wish.UserID != o.AuthedUser {
// 		return nil, &RequestError{
// 			Status: http.StatusUnauthorized,
// 			Err:    ErrUserNotAuthorized,
// 		}
// 	}

// 	db = lib.DB.Model(&wish).Updates(&Wish{
// 		Name:        b.Name,
// 		Description: b.Description,
// 		Link:        b.Link,
// 		Image:       b.Image,
// 	})
// 	if db.Error != nil {
// 		lib.LogError(lib.LPanic, "Could not update wish", db.Error)
// 	}

// 	return &Success{
// 		Status: http.StatusOK,
// 		View: &views.UWish{
// 			ID:          id,
// 			Name:        b.Name,
// 			Description: b.Description,
// 			Link:        b.Link,
// 			Image:       b.Image,
// 		},
// 	}, nil
// }

// // DeleteWish is used to delete a wish from the database
// func DeleteWish(o *Options) (*Success, error) {
// 	var wish Wish

// 	id := o.Params["id"].(uint64)

// 	db := lib.DB.Select("id, user_id").First(&wish, id)
// 	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
// 		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
// 	} else if db.RecordNotFound() {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrWishNotFound,
// 		}
// 	}
// 	if wish.UserID != o.AuthedUser {
// 		return nil, &RequestError{
// 			Status: http.StatusUnauthorized,
// 			Err:    ErrUserNotAuthorized,
// 		}
// 	}

// 	db = lib.DB.Delete(wish)
// 	if db.Error != nil {
// 		lib.LogError(lib.LPanic, "Could not delete wish", db.Error)
// 	}

// 	return &Success{
// 		Status: http.StatusOK,
// 	}, nil
// }

// func AddWantToFulfill(o *Options) (*Success, error) { // TODO: User should not be able to add itself to WantToFulfill when it is already a claimer or fulfiller
// 	var wish Wish

// 	id := o.Params["id"].(uint64)

// 	db := lib.DB.Select("id, user_id").First(&wish, id)
// 	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
// 		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
// 	} else if db.RecordNotFound() {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrWishNotFound,
// 		}
// 	}

// 	if o.AuthedUser == wish.UserID || !AreFriends(wish.UserID, o.AuthedUser) {
// 		return nil, &RequestError{
// 			Status: http.StatusUnauthorized,
// 			Err:    ErrUserNotAuthorized,
// 		}
// 	}

// 	asso := lib.DB.Model(&Wish{ID: id}).Where("user_id = ?", o.AuthedUser).Association("WantToFulfill")

// 	if asso.Count() != 0 {
// 		return nil, &RequestError{
// 			Status: http.StatusConflict,
// 			Err:    ErrUserExists,
// 		}
// 	}

// 	err := lib.DB.Model(&Wish{ID: id}).Association("WantToFulfill").Append(&User{ID: o.AuthedUser}).Error
// 	if err != nil {
// 		lib.LogError(lib.LPanic, "Could not add to WantToFulfill", err)
// 	}

// 	return &Success{
// 		Status: http.StatusOK,
// 		View: &views.WishID{
// 			ID: wish.ID,
// 		},
// 	}, nil
// }

// func AddClaimer(o *Options) (*Success, error) {
// 	id := o.Params["id"].(uint64)

// 	asso := lib.DB.Model(&Wish{ID: id}).Where("user_id = ?", o.AuthedUser).Association("WantToFulfill")
// 	if asso.Error != nil && !gorm.IsRecordNotFoundError(asso.Error) {
// 		lib.LogError(lib.LPanic, "Could not read wish's WantToFulfill", asso.Error)
// 	}

// 	if asso.Count() != 1 {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrUserNotFound,
// 		}
// 	}

// 	err := lib.DB.Transaction(func(tx *gorm.DB) error {
// 		asso := tx.Model(&Wish{ID: id}).Association("Claimers").Append(&User{ID: o.AuthedUser})
// 		if asso.Error != nil {
// 			return asso.Error
// 		}

// 		asso = tx.Model(&Wish{ID: id}).Association("WantToFulfill").Delete(&User{ID: o.AuthedUser})
// 		if asso.Error != nil {
// 			return asso.Error
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		lib.LogError(lib.LPanic, "Could not add to Claimers", err)
// 	}

// 	return &Success{
// 		Status: http.StatusOK,
// 		View: &views.WishID{
// 			ID: id,
// 		},
// 	}, nil
// }

// func AcceptClaimer(o *Options) (*Success, error) {
// 	var wish Wish

// 	id := o.Params["id"].(uint64)
// 	b := o.B.(*bindings.Claimer)

// 	db := lib.DB.Select("id, user_id").First(&wish, id)
// 	if db.RecordNotFound() {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrWishNotFound,
// 		}
// 	}

// 	if o.AuthedUser != wish.UserID {
// 		return nil, &RequestError{
// 			Status: http.StatusUnauthorized,
// 			Err:    ErrUserNotAuthorized,
// 		}
// 	}

// 	count := lib.DB.Model(&wish).Where("user_id = ?", b.ID).Association("Claimers").Count()
// 	if count != 1 {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrUserNotFound,
// 		}
// 	}

// 	err := lib.DB.Transaction(func(tx *gorm.DB) error {
// 		asso := tx.Model(&wish).Association("Fulfillers").Append(&User{ID: b.ID})
// 		if asso.Error != nil {
// 			return asso.Error
// 		}

// 		asso = tx.Model(&wish).Association("Claimers").Delete(&User{ID: b.ID})
// 		if asso.Error != nil {
// 			return asso.Error
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		lib.LogError(lib.LPanic, "Could not accept fulfillment claim", err)
// 	}

// 	return &Success{
// 		Status: http.StatusOK,
// 		View: &views.WishID{
// 			ID: id,
// 		},
// 	}, nil
// }

// func RejectClaimer(o *Options) (*Success, error) {
// 	var wish Wish

// 	id := o.Params["id"].(uint64)
// 	b := o.B.(*bindings.Claimer)

// 	db := lib.DB.Select("id, user_id").First(&wish, id)
// 	if db.RecordNotFound() {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrWishNotFound,
// 		}
// 	}

// 	if o.AuthedUser != wish.UserID {
// 		return nil, &RequestError{
// 			Status: http.StatusUnauthorized,
// 			Err:    ErrUserNotAuthorized,
// 		}
// 	}

// 	count := lib.DB.Model(&wish).Where("user_id = ?", b.ID).Association("Claimers").Count()
// 	if count != 1 {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrUserNotFound,
// 		}
// 	}

// 	err := lib.DB.Transaction(func(tx *gorm.DB) error {
// 		asso := tx.Model(&wish).Association("WantToFulfill").Append(&User{ID: b.ID})
// 		if asso.Error != nil {
// 			return asso.Error
// 		}

// 		asso = tx.Model(&wish).Association("Claimers").Delete(&User{ID: b.ID})
// 		if asso.Error != nil {
// 			return asso.Error
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		lib.LogError(lib.LPanic, "Could not reject fulfillment claim", err)
// 	}

// 	return &Success{
// 		Status: http.StatusOK,
// 		View: &views.WishID{
// 			ID: id,
// 		},
// 	}, nil
// }

// func ReadClaimers(o *Options) (*Success, error) {
// 	var wish Wish
// 	var claimers []User
// 	var vs []*views.RUser

// 	id := o.Params["id"].(uint64)

// 	db := lib.DB.Select("id, user_id").First(&wish, id)
// 	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
// 		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
// 	} else if db.RecordNotFound() {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrWishNotFound,
// 		}
// 	}

// 	if o.AuthedUser != wish.UserID {
// 		return nil, &RequestError{
// 			Status: http.StatusUnauthorized,
// 			Err:    ErrUserNotAuthorized,
// 		}
// 	}

// 	lib.DB.Model(&Wish{ID: id}).Association("Claimers").Find(&claimers)

// 	for _, c := range claimers {
// 		vs = append(vs, &views.RUser{
// 			ID:        c.ID,
// 			FirstName: c.FirstName,
// 			LastName:  c.LastName,
// 		})
// 	}

// 	return &Success{
// 		Status: http.StatusOK,
// 		View: &views.ReadUsers{
// 			Users: vs,
// 		},
// 	}, nil
// }

// func ReadFulfillers(o *Options) (*Success, error) {
// 	var wish Wish
// 	var fulfillers []User
// 	var vs []*views.RUser

// 	id := o.Params["id"].(uint64)

// 	db := lib.DB.Select("id, user_id").First(&wish, id)
// 	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
// 		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
// 	} else if db.RecordNotFound() {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrWishNotFound,
// 		}
// 	}

// 	if o.AuthedUser != wish.UserID {
// 		return nil, &RequestError{
// 			Status: http.StatusUnauthorized,
// 			Err:    ErrUserNotAuthorized,
// 		}
// 	}

// 	lib.DB.Model(&Wish{ID: id}).Association("Fulfillers").Find(&fulfillers)

// 	for _, c := range fulfillers {
// 		vs = append(vs, &views.RUser{
// 			ID:        c.ID,
// 			FirstName: c.FirstName,
// 			LastName:  c.LastName,
// 		})
// 	}

// 	return &Success{
// 		Status: http.StatusOK,
// 		View: &views.ReadUsers{
// 			Users: vs,
// 		},
// 	}, nil
// }

// func CountWantToFulfill(o *Options) (*Success, error) {
// 	var wish Wish

// 	id := o.Params["id"].(uint64)

// 	db := lib.DB.Select("id").First(&wish, id)
// 	if db.Error != nil && !gorm.IsRecordNotFoundError(db.Error) {
// 		lib.LogError(lib.LPanic, "Could not read wish", db.Error)
// 	} else if db.RecordNotFound() {
// 		return nil, &RequestError{
// 			Status: http.StatusNotFound,
// 			Err:    ErrWishNotFound,
// 		}
// 	}

// 	count := lib.DB.Model(&wish).Association("WantToFulfill").Count()

// 	return &Success{
// 		Status: http.StatusOK,
// 		View: &views.WishCount{
// 			Count: count,
// 		},
// 	}, nil
// }

func init() {
	db.DB.AutoMigrate(&Wish{})
}
