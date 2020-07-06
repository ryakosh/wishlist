package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/bindings"
	"github.com/ryakosh/wishlist/lib/models"
)

// ErrRequestIsInvalid is returned when the provided request could not
// be handled due to validation or parsing errors
var ErrRequestIsInvalid = errors.New("Request is invalid")

var r *gin.Engine

func genHandler(modelFunc interface{}, b interface{}, params []string, needAuthedUser bool) gin.HandlerFunc {
	o := new(models.Options)

	return func(c *gin.Context) {
		if needAuthedUser {
			authedUser := authedUser(c)
			if authedUser == "" {
				return
			}

			o.AuthedUser = authedUser
		}

		if b != nil {
			if ok := bindJSON(c, &b); !ok {
				return
			}

			o.B = b
		}

		if params != nil && len(params) != 0 {
			cp, err := canonicalizeParams(params, c)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": ErrRequestIsInvalid.Error(),
				})
				return
			}

			o.Params = cp
		}

		mf := modelFunc.(func(*models.Options) (*models.Success, error))

		view, err := mf(o)
		if err != nil {
			err := err.(*models.RequestError)
			c.JSON(err.Status, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(view.Status, view.View)
	}
}

func canonicalizeParams(params []string, c *gin.Context) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	for _, p := range params {
		s := strings.SplitN(p, ":", 2)

		if len(s) == 1 || s[1] == "string" {
			m[s[0]] = c.Param(s[0])
		} else if s[1] == "uint64" {
			u, err := strconv.ParseUint(c.Param(s[0]), 10, 64)
			if err != nil {
				return nil, err
			}

			m[s[0]] = u
		}
	}

	return m, nil
}

func bindJSON(c *gin.Context, bindTo interface{}) bool {
	if err := c.ShouldBindJSON(bindTo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrRequestIsInvalid.Error(),
		})

		return false
	}

	return true
}

// authedUser returns the authenticated user or empty string otherwise
func authedUser(c *gin.Context) string {
	authedUser, ok := c.Get(models.UserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": models.ErrUserNotAuthorized.Error(),
		})

		return ""
	}

	return authedUser.(string)
}

func main() {
	r.Run(":8080")
}

func init() {
	serverLog, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		lib.LogError(lib.LFatal, "Could not create log file", err)

	}
	log.SetOutput(serverLog)

	gin.DisableConsoleColor()
	ginLog, err := os.OpenFile("gin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		lib.LogError(lib.LFatal, "Could not create log file", err)
	}
	gin.DefaultWriter = ginLog

	r = gin.Default()
	r.Use(AccessLogger())

	r.POST("/login", genHandler(models.LoginUser, new(bindings.LoginUser), nil, false))
	users := r.Group("/users")
	users.POST("", genHandler(models.CreateUser, new(bindings.CUser), nil, false))
	users.GET("/:id", genHandler(models.ReadUser, nil, []string{"id"}, false))
	users.PUT("/:id", models.Authenticate(), genHandler(models.UpdateUser, new(bindings.UUser), []string{"id"}, true))
	users.DELETE("/:id", models.Authenticate(), genHandler(models.DeleteUser, nil, []string{"id"}, true))
	users.PUT("/:id/verify_email", models.Authenticate(), genHandler(models.VerifyUserEmail, new(bindings.VerifyUserEmail), []string{"id"}, true))

	friendRequests := users.Group("/:id/friend_requests")
	friendRequests.GET("", models.Authenticate(), models.RequireEmailVerification(), genHandler(models.ReadFriendRequests, nil, []string{"id"}, true))
	friendRequests.GET("/count", models.Authenticate(), models.RequireEmailVerification(), genHandler(models.CountFriendRequests, nil, []string{"id"}, true))

	friends := users.Group("/:id/friends")
	friends.GET("", models.Authenticate(), models.RequireEmailVerification(), genHandler(models.ReadFriends, nil, []string{"id"}, true))
	friends.GET("/count", models.Authenticate(), models.RequireEmailVerification(), genHandler(models.CountFriends, nil, []string{"id"}, true))
	friends.PUT("/send_request", models.Authenticate(), models.RequireEmailVerification(), genHandler(models.ReqFriendship, new(bindings.Requestee), []string{"id"}, true))
	friends.DELETE("/undo_request", models.Authenticate(), models.RequireEmailVerification(), genHandler(models.UnReqFriendship, new(bindings.Requestee), []string{"id"}, true))
	friends.PUT("/accept_request", models.Authenticate(), models.RequireEmailVerification(), genHandler(models.AccFriendship, new(bindings.Requestee), []string{"id"}, true))
	friends.DELETE("/reject_request", models.Authenticate(), models.RequireEmailVerification(), genHandler(models.RejFriendship, new(bindings.Requestee), []string{"id"}, true))

	wishes := r.Group("/wishes")
	wishes.POST("", models.Authenticate(), models.RequireEmailVerification(), genHandler(models.CreateWish, new(bindings.CWish), nil, true))
	wishes.GET("/:id", genHandler(models.ReadWish, nil, []string{"id:uint64"}, false))
	wishes.PUT("/:id", models.Authenticate(), genHandler(models.UpdateWish, new(bindings.UWish), []string{"id:uint64"}, true))
	wishes.DELETE("/:id", models.Authenticate(), genHandler(models.DeleteWish, nil, []string{"id:uint64"}, true))
	wishes.PUT("/:id/add_fulfiller", models.Authenticate(), models.RequireEmailVerification(), genHandler(models.AddWantToFulfill, nil, []string{"id:uint64"}, true))
	wishes.PUT("/:id/add_claimer", models.Authenticate(), genHandler(models.AddClaimer, nil, []string{"id:uint64"}, true))
	wishes.PUT("/:id/accept_claimer", models.Authenticate(), genHandler(models.AcceptClaimer, new(bindings.Claimer), []string{"id:uint64"}, true))
	wishes.PUT("/:id/reject_claimer", models.Authenticate(), genHandler(models.RejectClaimer, new(bindings.Claimer), []string{"id:uint64"}, true))
	wishes.GET("/:id/read_fulfillers", models.Authenticate(), genHandler(models.ReadFulfillers, nil, []string{"id:uint64"}, true))
	wishes.GET("/:id/read_claimers", models.Authenticate(), genHandler(models.ReadClaimers, nil, []string{"id:uint64"}, true))
	wishes.GET("/:id/count_fulfillers", genHandler(models.CountWantToFulfill, nil, []string{"id:uint64"}, false))
}
