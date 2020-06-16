package models

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/bindings"
	"github.com/ryakosh/wishlist/lib/views"
)

const (
	// TokenCookieKey is a cookie key used for authentication
	TokenCookieKey = "token"

	// UserKey is a gin context keystore key used to indicate
	// that a user is authenticated
	UserKey = "user"
)

var (
	// ErrUserExists is returned when User already exists in the database
	ErrUserExists = errors.New("User already exists")

	// ErrUserNotFound is returned when User does not exist in the database
	ErrUserNotFound = errors.New("User not found")

	// ErrUserNotAuthorized is returned when User is not authorized to access a resource
	ErrUserNotAuthorized = errors.New("User not authorized")

	// ErrUnmOrPwdIncorrect is returned when the provided username or
	// password are incorrect
	ErrUnmOrPwdIncorrect = errors.New("Username or password is incorrect")
)

var argonConfig = &argon2id.Params{
	Memory:      32 * 1024,
	Iterations:  2,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// User represents a user in the app
type User struct {
	ID                string `gorm:"type:varchar(64)"`
	Email             string `gorm:"varchar(254);unique"`
	IsEmailVerified   bool
	Password          string `gorm:"varchar(256)"`
	FirstName         string `gorm:"type:varchar(64)"`
	LastName          string `gorm:"type:varchar(64)"`
	Wishes            []Wish
	FulfilledWishes   []Wish `gorm:"foreignkey:FulfilledBy"`
	WantFulfillWishes []Wish `gorm:"many2many:userswant_wishes"`
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}

// CreateUser is used to register/add a user to the database
func CreateUser(b *bindings.CUser) (*views.CUser, error) {
	var user User

	db := lib.DB.Where("id = ?", b.ID).Or("email = ?", b.Email).First(&user)
	if db.RecordNotFound() {
		user = User{
			ID:        b.ID,
			Email:     b.Email,
			Password:  genPasswordHash(b.Password),
			FirstName: b.FirstName,
			LastName:  b.LastName,
		}

		lib.DB.Create(&user)
		return &views.CUser{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}, nil
	}

	return nil, ErrUserExists
}

// ReadUser is used to get general information about a user in the database
func ReadUser(id string) (*views.RUser, error) {
	var user User

	db := lib.DB.Select("id, first_name, last_name").Where("id = ?", id).First(&user)
	if !db.RecordNotFound() {
		return &views.RUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}, nil
	}

	return nil, ErrUserNotFound
}

// UpdateUser is used to update user's general information
func UpdateUser(b *bindings.UUser, authedUser string) *views.UUser {
	user := User{
		ID: authedUser,
	}

	lib.DB.Model(&user).Updates(&User{
		FirstName: b.FirstName,
		LastName:  b.LastName,
	})

	return &views.UUser{
		FirstName: b.FirstName,
		LastName:  b.LastName,
	}
}

// DeleteUser is used to delete a user from the database
func DeleteUser(authedUser string) {
	user := &User{
		ID: authedUser,
	}

	lib.DB.Delete(user)
}

// LoginUser is used for user authentication
func LoginUser(b *bindings.LoginUser) (string, error) {
	var user User

	db := lib.DB.Select("id, password").Where("id = ?", b.ID).First(&user)
	if !db.RecordNotFound() && verifyPassword(b.Password, user.Password) {
		return lib.Encode(user.ID), nil
	}

	return "", ErrUnmOrPwdIncorrect
}

func genPasswordHash(password string) string {
	hash, err := argon2id.CreateHash(password, argonConfig)
	if err != nil {
		log.Panicf("error: Could not generate password's hash\n\treason: %s", err)
	}

	return hash
}

func verifyPassword(password string, hash string) bool {
	isMatch, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Panicf("error: Could not verify password\n\treason: %s", err)
	}

	return isMatch
}

// Authenticate is a middleware that is used to authenticate users
// on certain endpoints using cookies, it's not enforcing authentication
// on endpoints that it's beeing used so endpoints should decide whether
// they require authentication or not, however it aborts requests if
// the provided token is malformed, expired or not valid
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(TokenCookieKey)
		if err != nil {
			c.Next()
			return
		}

		claims, valid, err := lib.Decode(token)
		if err == nil && valid {
			var user User
			sub := claims["sub"]

			db := lib.DB.Select("id").Where("id = ?", sub).First(&user)
			if !db.RecordNotFound() {
				c.Set(UserKey, sub)
				c.Next()
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrUserNotFound.Error(),
			})
		} else if lib.IsMalformed(err) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": lib.ErrTokenIsMalformed.Error(),
			})
		} else if lib.HasExpired(err) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": lib.ErrTokenHasExpired.Error(),
			})
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": lib.ErrTokenHasExpired.Error(),
			})
		}
	}
}

func init() {
	lib.DB.AutoMigrate(&User{})
}
