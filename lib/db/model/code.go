package model

import (
	"errors"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/db"
)

const (
	// CodeMaxRetries is used to set the maximum times a user can enter
	// their code without the code getting expired
	CodeMaxRetries = 3

	// CodeTTL is used to set code's Time-To-Live, after this duration
	// the code is not valid anymore
	CodeTTL = time.Minute * 30
)

var (
	// ErrCodeExists is returned when code already exists in the database
	ErrCodeExists = errors.New("Code already exists")

	// ErrCodeNotFound is returned when code does not exist in the database
	ErrCodeNotFound = errors.New("Code not found")

	// ErrCodeNotMatch is returned when the provided code by user does not
	// match the one in the database
	ErrCodeNotMatch = errors.New("Code does not match")
)

// Code is a table that stores safe random codes that are
// used for verifying emails, or when users forget their password .etc
type Code struct {
	UserID     string `gorm:"primary_key"`
	Code       string
	RetryCount uint
	CreatedAt  *time.Time
}

// CreateCode is used to create a new safe random code in the database
func CreateCode(username string) (*Success, error) {
	var user User
	var code Code

	d := db.DB.Select("id").Where("id = ?", username).First(&user)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	} else if d.RecordNotFound() {
		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrUserNotFound,
		}
	}

	d = db.DB.Select("user_id, created_at").Where("user_id = ?", username).First(&code)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read code", d.Error)
	}

	if !d.RecordNotFound() {
		now := time.Now().UTC()
		deadline := code.CreatedAt.UTC().Add(CodeTTL)

		if !now.After(deadline) {
			return nil, &RequestError{
				Status: http.StatusConflict,
				Err:    ErrCodeExists,
			}
		}

		d := db.DB.Delete(&code)
		if d.Error != nil {
			lib.LogError(lib.LPanic, "Could not delete code", d.Error)
		}
	}

	randCode, err := lib.GenRandCode(10)
	if err != nil {
		return nil, &ServerError{
			Status: http.StatusInternalServerError,
			Reason: err,
		}
	}

	code = Code{
		UserID: username,
		Code:   randCode,
	}

	d = db.DB.Create(&code)
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not create code", d.Error)
	}

	return &Success{
		Status: http.StatusCreated,
		View:   code.Code,
	}, nil
}

// VerifyCode is used to compare the provided random code by user
// with the random code in the database
func VerifyCode(username string, randCode string) (*Success, error) {
	var code Code

	d := db.DB.Select("user_id, code, retry_count, created_at").Where("user_id = ?", username).First(&code)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read code", d.Error)
	} else if d.RecordNotFound() {
		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrCodeNotFound,
		}
	}

	now := time.Now().UTC()
	deadline := code.CreatedAt.UTC().Add(CodeTTL)
	if code.RetryCount == CodeMaxRetries || now.After(deadline) {
		d := db.DB.Delete(&code)
		if d.Error != nil {
			lib.LogError(lib.LPanic, "Could not delete code", d.Error)
		}

		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrCodeNotFound,
		}
	}

	if code.Code != randCode {
		rc := code.RetryCount + 1
		d := db.DB.Model(&code).Update("retry_count", rc)
		if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
			lib.LogError(lib.LPanic, "Could not update code", d.Error)
		}

		return nil, &RequestError{
			Status: http.StatusOK,
			Err:    ErrCodeNotMatch,
		}
	}

	d = db.DB.Delete(&code)
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not delete code", d.Error)
	}

	return &Success{
		Status: http.StatusOK,
		View:   true,
	}, nil
}

func init() {
	db.DB.AutoMigrate(&Code{})
}
