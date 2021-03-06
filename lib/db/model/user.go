package model

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/alexedwards/argon2id"
	"github.com/jinzhu/gorm"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/db"
)

type key int

const (
	authedUserKey key = iota
)

const (
	UserWishesAsso         db.Association = "Wishes"
	UserCodeAsso           db.Association = "Code"
	UserFriendsAsso        db.Association = "Friends"
	UserFriendRequestsAsso db.Association = "FriendRequests"
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

	// ErrBearerTokenMalformed is returned when the provided bearer token in
	// Authorization header is malformed
	ErrBearerTokenMalformed = errors.New("Bearer token is malformed")

	// ErrUserNotVerified is returned when user's email address has not yet
	// been verified
	ErrUserNotVerified = errors.New("User not verified")

	// ErrEmailVerified is returned when user's email address is
	// already verfied
	ErrEmailVerified = errors.New("Email is already verified")
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
	ID              string `gorm:"type:varchar(64)"`
	Email           string `gorm:"type:varchar(254);unique"`
	IsEmailVerified bool
	Password        string  `gorm:"type:varchar(256)"`
	FirstName       *string `gorm:"type:varchar(64)"`
	LastName        *string `gorm:"type:varchar(64)"`
	Wishes          []Wish  `gorm:"foreignkey:Owner"`
	Code            Code
	Friends         []*User `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id"`
	FriendRequests  []*User `gorm:"many2many:friendrequests;association_jointable_foreignkey:requester_id"`
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

// AfterDelete is used to clean up after the user got deleted
func (u *User) AfterDelete(tx *gorm.DB) error {
	d := db.DB.Where("user_id = ?", u.ID).Delete(&Wish{})
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not delete user's wishes", d.Error)
	}

	d = db.DB.Where("user_id = ?", u.ID).Delete(&Code{})
	if d.Error != nil {
		lib.LogError(lib.LPanic, "Could not delete user's code", d.Error)

	}

	return nil
}

func GenPasswordHash(password string) string {
	hash, err := argon2id.CreateHash(password, argonConfig)
	if err != nil {
		lib.LogError(lib.LPanic, "Could not generate password's hash", err)
	}

	return hash
}

func VerifyPassword(password string, hash string) bool {
	isMatch, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		lib.LogError(lib.LPanic, "Could not verify password", err)
	}

	return isMatch
}

func AreFriends(user, authedUser string) bool {
	count := db.DB.Model(&User{ID: authedUser}).Where("friend_id = ?", user).Association("Friends").Count()
	if count != 1 {
		return false
	}

	return true
}

// AuthRequired is a middleware that is used to authenticate users
// on certain endpoints using Authorization header, it's not enforcing authentication
// on endpoints that it's beeing used so endpoints should decide whether
// they require authentication or not, however it aborts requests if
// the provided token is malformed, expired or not valid
func AuthRequired(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	authedUser := AuthedUserFromCtx(ctx)
	if authedUser != "" {
		return authedUser, nil
	}

	c := lib.GinCtxFromCtx(ctx)
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		return nil, ErrUserNotAuthorized
	}

	token := strings.Fields(authorizationHeader)
	if len(token) != 2 || token[0] != "Bearer" {
		return nil, ErrBearerTokenMalformed
	}

	claims, valid, err := lib.Decode(token[1])
	if err == nil && valid {
		var user User
		sub := claims["sub"]

		d := db.DB.Select("id").Where("id = ? AND email = ?", sub, claims["email"]).First(&user)
		if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
			lib.LogError(lib.LPanic, "Could not read user", d.Error)
		} else if d.RecordNotFound() {
			return nil, ErrUserNotFound
		}

		authedUser = sub.(string)
		ctx = context.WithValue(ctx, authedUserKey, authedUser)

		return next(ctx)
	} else if lib.IsMalformed(err) {
		return nil, lib.ErrTokenIsMalformed
	} else if lib.HasExpired(err) {
		return nil, lib.ErrTokenHasExpired
	} else {
		return nil, lib.ErrTokenIsInvalid
	}
}

func AuthedUserFromCtx(ctx context.Context) string {
	authedUser := ctx.Value(authedUserKey)
	if authedUser == nil {
		return ""
	}

	c, ok := authedUser.(string)
	if !ok {
		lib.LogError(lib.LFatal, "AuthedUser has wrong type", nil)
	}

	return c
}

// EmailVerificationRequired is a middleware that is used to
// indicate that a user's email address must be verified in order
// to access this endpoint, it should be called after the Authenticate
// middleware
func EmailVerificationRequired(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	var user User
	authedUser := AuthedUserFromCtx(ctx)

	d := db.DB.Select("is_email_verified").Where("id = ?", authedUser).First(&user)
	if d.Error != nil && !gorm.IsRecordNotFoundError(d.Error) {
		lib.LogError(lib.LPanic, "Could not read user", d.Error)
	}

	if !user.IsEmailVerified {
		return nil, ErrUserNotAuthorized
	}

	return next(ctx)
}

func init() {
	db.DB.AutoMigrate(&User{})
}
