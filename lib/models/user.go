package models

import (
	"errors"
	"fmt"
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
	ErrUserExists = errors.New("error: User already exists")

	// ErrUserNotFound is returned when User does not exist in the database
	ErrUserNotFound = errors.New("error: User not found")

	// ErrUserNotAuthorized is returned when User is not authorized to access a resource
	ErrUserNotAuthorized = errors.New("error: User not authorized")
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

// CreateUser is used to add a User to our database
func CreateUser(b *bindings.CuUser) (*views.CuUser, error) {
	var user *User

	lib.DB.Where("id = ?", b.ID).Or("email = ?", b.Email).First(&user)
	if user == nil {
		user = &User{
			ID:        b.ID,
			Email:     b.Email,
			Password:  genPasswordHash(b.Password),
			FirstName: b.FirstName,
			LastName:  b.LastName,
		}

		lib.DB.Create(&user)
		return &views.CuUser{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}, nil
	}

	return nil, ErrUserExists
}

// ReadUser is used to get information about a User in our database
func ReadUser(b *bindings.RUser) (*views.RUser, error) {
	var user *User

	lib.DB.Select("id", "first_name", "last_name").Where("id = ?", b.ID).First(&user)
	if user != nil {
		return &views.RUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}, nil
	}

	return nil, ErrUserNotFound
}

// UpdateUser is used to update a User in our database
func UpdateUser(b *bindings.CuUser, authedUser string) *views.CuUser {
	user := &User{
		ID: authedUser,
	}

	lib.DB.Model(&user).Select("password,first_name,last_name").Updates(&User{
		Password:  genPasswordHash(b.Password),
		FirstName: b.FirstName,
		LastName:  b.LastName,
	})

	return &views.CuUser{
		ID:        authedUser,
		FirstName: b.FirstName,
		LastName:  b.LastName,
	}
}

// DeleteUser is used to delete a User from our database
func DeleteUser(authedUser string) {
	user := &User{
		ID: authedUser,
	}

	lib.DB.Delete(user)
}

func genPasswordHash(password string) string {
	hash, err := argon2id.CreateHash(password, argonConfig)
	if err != nil {
		panic(fmt.Sprintf("error: Could not generate password's hash\n\treason: %s", err))
	}

	return hash
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
		} else {
			claims, valid, err := lib.Decode(token)
			if err == nil && valid {
				var user *User
				sub := (*claims)["sub"]

				lib.DB.Where("id = ?", sub).First(user)
				if user != nil {
					c.Set(UserKey, sub)
					c.Next()
				}
			} else if lib.IsMalformed(err) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token malformed",
				})
			} else if lib.HasExpired(err) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token expired",
				})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token is invalid",
				})
			}
		}
	}
}

func init() {
	lib.DB.AutoMigrate(&User{})
}
