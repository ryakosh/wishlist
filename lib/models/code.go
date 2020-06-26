package models

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ryakosh/wishlist/lib"
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
func CreateCode(username string) error {
	var user User
	var code Code

	db := lib.DB.Select("id").Where("id = ?", username).First(&user)
	if db.RecordNotFound() {
		return &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrUserNotFound,
		}
	}

	db = lib.DB.Select("user_id, created_at").Where("user_id = ?", username).First(&code)
	if !db.RecordNotFound() {
		now := time.Now().UTC()
		deadline := code.CreatedAt.UTC().Add(CodeTTL)

		if !now.After(deadline) {
			return &RequestError{
				Status: http.StatusConflict,
				Err:    ErrCodeExists,
			}
		}

		lib.DB.Delete(&code)
	}

	randCode, err := lib.GenRandCode(20)
	if err != nil {
		log.Panicf("error: Could not generate safe random code\n\treason: %s", err)
	}

	code = Code{
		UserID: username,
		Code:   randCode,
	}

	lib.DB.Create(&code)

	return nil
}

// VerifyCode is used to compare the provided random code by user
// with the random code in the database
func VerifyCode(username string, randCode string) (bool, error) {
	var code Code

	db := lib.DB.Select("user_id, code, retry_count").Where("user_id = ?", username).First(&code)
	if db.RecordNotFound() {
		return false, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrCodeNotFound,
		}
	}

	now := time.Now().UTC()
	deadline := code.CreatedAt.UTC().Add(CodeTTL)
	if code.RetryCount == CodeMaxRetries || now.After(deadline) {
		lib.DB.Delete(&code)

		return false, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrCodeNotFound,
		}
	}

	if code.Code != randCode {
		rc := code.RetryCount + 1
		lib.DB.Model(&code).Update("retry_count", rc)

		return false, &RequestError{
			Status: http.StatusBadRequest,
			Err:    ErrCodeNotMatch,
		}
	}

	return true, nil
}

func init() {
	lib.DB.AutoMigrate(&Code{})
}
