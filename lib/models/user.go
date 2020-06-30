package models

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/bindings"
	"github.com/ryakosh/wishlist/lib/email"
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
	ID                string `gorm:"type:varchar(64)"`
	Email             string `gorm:"varchar(254);unique"`
	IsEmailVerified   bool
	Password          string `gorm:"varchar(256)"`
	FirstName         string `gorm:"type:varchar(64)"`
	LastName          string `gorm:"type:varchar(64)"`
	Wishes            []Wish
	FulfilledWishes   []Wish `gorm:"foreignkey:FulfilledBy"`
	WantFulfillWishes []Wish `gorm:"many2many:userswant_wishes"`
	Code              Code
	Friends           []*User `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id"`
	FriendRequests    []*User `gorm:"many2many:friendrequests;association_jointable_foreignkey:requester_id"`
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}

// AfterDelete is used to clean up after the user got deleted
func (u *User) AfterDelete(tx *gorm.DB) error {
	lib.DB.Where("user_id = ?", u.ID).Delete(&Wish{})
	lib.DB.Where("user_id = ?", u.ID).Delete(&Code{})

	return nil
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

		code, err := CreateCode(user.ID)
		if err != nil {
			se, ok := err.(*ServerError)
			if ok {
				log.Printf("error: Could not generate email confirmation mail\n\treason: %s\n", err)
				return nil, &RequestError{
					Status: se.Status,
					Err:    email.ErrSendMail,
				}
			}

			return nil, err
		}

		mail, err := email.GenEmailConfirmMail(user.ID, code)
		if err != nil {
			log.Printf("error: Could not generate email confirmation mail\n\treason: %s\n", err)
			return nil, &RequestError{
				Status: http.StatusInternalServerError,
				Err:    email.ErrSendMail,
			}
		}

		err = email.Send(email.BotEmailEnv, user.Email, "لطفا ایمیل خود را تایید کنید [ویش لیست]", mail)
		if err != nil {
			log.Printf("error: Could not generate email confirmation mail\n\treason: %s\n", err)
			return nil, &RequestError{
				Status: http.StatusInternalServerError,
				Err:    email.ErrSendMail,
			}
		}

		return &views.CUser{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}, nil
	}

	return nil, &RequestError{
		Status: http.StatusConflict,
		Err:    ErrUserExists,
	}
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

	return nil, &RequestError{
		Status: http.StatusNotFound,
		Err:    ErrUserNotFound,
	}
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

	db := lib.DB.Select("id, email, password, is_email_verified").Where("id = ?", b.ID).First(&user)
	if !db.RecordNotFound() && verifyPassword(b.Password, user.Password) {
		return lib.Encode(user.ID, user.Email), nil
	}

	return "", &RequestError{
		Status: http.StatusNotFound,
		Err:    ErrUnmOrPwdIncorrect,
	}
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

// VerifyUserEmail is used to verify user's email address using
// the generated safe random code
func VerifyUserEmail(b *bindings.VerifyUserEmail, authedUser string) error {
	var user User

	lib.DB.Select("is_email_verified").Where("id = ?", authedUser).First(&user)

	if user.IsEmailVerified {
		return &RequestError{
			Status: http.StatusOK,
			Err:    ErrEmailVerified,
		}
	}

	isMatch, err := VerifyCode(authedUser, b.Code)
	if err != nil {
		return err
	}

	if isMatch {
		lib.DB.Model(&User{ID: authedUser}).Update("is_email_verified", true)
	}

	return nil
}

// ReqFriendship is used to request friendship from another user in the
// database
func ReqFriendship(b *bindings.Requestee, authedUser string) (*views.Requestee, error) {
	var requestee User
	var friendsCount uint8
	var friendRequestsCount uint8

	if authedUser == b.Requestee {
		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrUserNotFound,
		}
	}

	db := lib.DB.Select("id").Where("id = ?", b.Requestee).First(&requestee)
	if db.RecordNotFound() {
		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrUserNotFound,
		}
	}

	lib.DB.Table("friendrequests").Where("user_id = ? AND requester_id = ?", requestee.ID, authedUser).Count(&friendRequestsCount)
	lib.DB.Table("friendships").Where("user_id = ? AND friend_id = ?", authedUser, requestee.ID).Count(&friendsCount)
	if friendRequestsCount != 0 || friendsCount != 0 {
		return nil, &RequestError{
			Status: http.StatusConflict,
			Err:    ErrUserExists,
		}
	}

	err := lib.DB.Model(&User{ID: requestee.ID}).Association("FriendRequests").Append(&User{ID: authedUser}).Error
	if err != nil {
		log.Printf("error: Could not request friendship\n\treason: %s", err)
		return nil, &RequestError{
			Status: http.StatusInternalServerError,
			Err:    ErrInternalServer,
		}
	}

	return &views.Requestee{
		Requestee: requestee.ID,
	}, nil
}

// AccFriendship is used to accept a friendship request from another user
// that has been previously requested for friendship
func AccFriendship(b *bindings.Requestee, authedUser string) (*views.Requestee, error) {
	var requestees []User

	lib.DB.Model(&User{ID: authedUser}).Select("id").Where(
		"requester_id = ?", b.Requestee).Related(&requestees, "FriendRequests")

	if len(requestees) != 1 {
		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrUserNotFound,
		}
	}

	err := lib.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&User{ID: authedUser}).Association("Friends").Append(requestees[0]).Error
		if err != nil {
			return err
		}

		err = tx.Model(&User{ID: requestees[0].ID}).Association("Friends").Append(&User{ID: authedUser}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&User{ID: authedUser}).Association("FriendRequests").Delete(requestees[0]).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("error: Could not accept friendship\n\treason: %s", err)
		return nil, &RequestError{
			Status: http.StatusInternalServerError,
			Err:    ErrInternalServer,
		}
	}

	return &views.Requestee{
		Requestee: requestees[0].ID,
	}, nil
}

// RejFriendship is used to reject a friendship request from another user
// that has been previously requested for friendship
func RejFriendship(b *bindings.Requestee, authedUser string) (*views.Requestee, error) {
	var requestees []User

	lib.DB.Model(&User{ID: authedUser}).Select("id").Where(
		"requester_id = ?", b.Requestee).Related(&requestees, "FriendRequests")

	if len(requestees) != 1 {
		return nil, &RequestError{
			Status: http.StatusNotFound,
			Err:    ErrUserNotFound,
		}
	}
	err := lib.DB.Model(&User{ID: authedUser}).Association("FriendRequests").Delete(requestees[0]).Error

	if err != nil {
		log.Printf("error: Could not reject friendship\n\treason: %s", err)
		return nil, &RequestError{
			Status: http.StatusInternalServerError,
			Err:    ErrInternalServer,
		}
	}

	return &views.Requestee{
		Requestee: requestees[0].ID,
	}, nil
}

// CountFriendRequests is used to count user's friend requests
func CountFriendRequests(authedUser string) *views.CountFriends {
	count := lib.DB.Model(&User{ID: authedUser}).Association("FriendRequests").Count()

	return &views.CountFriends{
		Count: count,
	}
}

// CountFriends is used to count user's friends
func CountFriends(authedUser string) *views.CountFriends {
	count := lib.DB.Model(&User{ID: authedUser}).Association("Friends").Count()

	return &views.CountFriends{
		Count: count,
	}
}

// ReadFriends is used to get user's friends
func ReadFriends(page uint64, authedUser string) *views.ReadFriends {
	var friends []User
	var vs []*views.RUser

	lib.DB.Model(&User{ID: authedUser}).Select(
		"id, first_name, last_name").Offset(
		(page * 10) - 10).Limit(10).Association("Friends").Find(&friends)

	for _, u := range friends {
		vs = append(vs, &views.RUser{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		})
	}

	return &views.ReadFriends{
		Friends: vs,
	}
}

// ReadFriendRequests is used to get user's friend requests
func ReadFriendRequests(page uint64, authedUser string) *views.ReadFriendRequests {
	var reqs []User
	var vs []*views.RUser

	lib.DB.Model(&User{ID: authedUser}).Select(
		"id, first_name, last_name").Offset(
		(page * 10) - 10).Limit(10).Association("FriendRequests").Find(&reqs)

	for _, u := range reqs {
		vs = append(vs, &views.RUser{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		})
	}

	return &views.ReadFriendRequests{
		Requesters: vs,
	}
}

// Authenticate is a middleware that is used to authenticate users
// on certain endpoints using Authorization header, it's not enforcing authentication
// on endpoints that it's beeing used so endpoints should decide whether
// they require authentication or not, however it aborts requests if
// the provided token is malformed, expired or not valid
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.Fields(c.GetHeader("Authorization"))

		if len(token) == 0 {
			c.Next()
			return
		} else if len(token) != 2 || token[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrBearerTokenMalformed.Error(),
			})
			return
		}

		claims, valid, err := lib.Decode(token[1])
		if err == nil && valid {
			var user User
			sub := claims["sub"]

			db := lib.DB.Select("id").Where("id = ? AND email = ?", sub, claims["email"]).First(&user)
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
				"error": lib.ErrTokenIsInvalid.Error(),
			})
		}
	}
}

// RequireEmailVerification is a middleware that is used to
// indicate that a user's email address must be verified in order
// to access this endpoint, it should be called after the Authenticate
// middleware
func RequireEmailVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		authedUser, ok := c.Get(UserKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrUserNotAuthorized.Error(),
			})
			return
		}

		var user User
		lib.DB.Select("is_email_verified").Where("id = ?", authedUser).First(&user)

		if !user.IsEmailVerified {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrUserNotAuthorized.Error(),
			})
		}
	}
}

func init() {
	lib.DB.AutoMigrate(&User{})
}
