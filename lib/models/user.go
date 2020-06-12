package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/bindings"
)

var (
	// ErrUserExists is returned when User already exists in the database
	ErrUserExists = errors.New("error: User already exists")

	// ErrUserNotFound is returned when User does not exist in the database
	ErrUserNotFound = errors.New("error: User not found")
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
	ID                string     `json:"username" gorm:"type:varchar(64)"`
	Email             string     `json:"email" gorm:"varchar(254);unique"`
	IsEmailVerified   bool       `json:"is_email_verfied"`
	Password          string     `json:"password" gorm:"varchar(256)"`
	FirstName         string     `json:"first_name" gorm:"type:varchar(64)"`
	LastName          string     `json:"last_name" gorm:"type:varchar(64)"`
	Wishes            []Wish     `json:"wishes"`
	FulfilledWishes   []Wish     `json:"fulfilled_wishes" gorm:"foreignkey:FulfilledBy"`
	WantFulfillWishes []Wish     `json:"want_fulfill_wishes" gorm:"many2many:userswant_wishes"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

// CreateUser is used to add a User to our database
func CreateUser(b *bindings.CuUser) (*lib.CuUserView, error) {
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
		return &lib.CuUserView{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}, nil
	}

	return nil, ErrUserExists
}

// ReadUser is used to get information about a User in our database
func ReadUser(b *bindings.RUser) (*lib.RdUserView, error) {
	var user *User

	lib.DB.Select("id", "first_name", "last_name").Where("id = ?", b.ID).First(&user)
	if user != nil {
		return &lib.RdUserView{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}, nil
	}

	return nil, ErrUserNotFound
}

// UpdateUser is used to update a User in our database
func UpdateUser(b *bindings.CuUser, authedUser string) *lib.CuUserView {
	user := &User{
		ID: authedUser,
	}

	lib.DB.Model(&user).Select("password,first_name,last_name").Updates(&User{
		Password:  genPasswordHash(b.Password),
		FirstName: b.FirstName,
		LastName:  b.LastName,
	})

	return &lib.CuUserView{
		ID:        authedUser,
		FirstName: b.FirstName,
		LastName:  b.LastName,
	}
}

// DeleteUser is used to delete a User from our database
func DeleteUser(authedUser string) *lib.RdUserView {
	user := &User{
		ID: authedUser,
	}

	lib.DB.Delete(user)

	return &lib.RdUserView{
		ID: authedUser,
	}
}

func genPasswordHash(password string) string {
	hash, err := argon2id.CreateHash(password, argonConfig)
	if err != nil {
		panic(fmt.Sprintf("error: Could not generate password's hash\n\treason: %s", err))
	}

	return hash
}

func init() {
	lib.DB.AutoMigrate(&User{})
}
