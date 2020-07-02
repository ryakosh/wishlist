package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib/bindings"
	"github.com/ryakosh/wishlist/lib/models"
)

// LoginUser is a route handler that is used for user authentication,
// it does so by responding with a json that contains the token
func LoginUser(c *gin.Context) { // TODO: Don't forget about CSRF attacks
	var b bindings.LoginUser

	if ok := bindJSON(c, &b); !ok {
		return
	}

	token, err := models.LoginUser(&b)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}

	// TODO: User should not be able to authenticate when they are authenticated already
	c.JSON(token.Status, gin.H{
		"token": token.View,
	})
}

// CreateUser is a route handler that is used to create/register a new user
func CreateUser(c *gin.Context) {
	var b bindings.CUser

	if ok := bindJSON(c, &b); !ok {
		return
	}

	view, err := models.CreateUser(&b)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}
	c.JSON(view.Status, view.View)
}

// ReadUser is a route handler that is used to get general information about a user
func ReadUser(c *gin.Context) {
	id := c.Param("id")

	view, err := models.ReadUser(id)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}
	c.JSON(view.Status, view.View)
}

// UpdateUser is a route handler that is used to update general information about a user
func UpdateUser(c *gin.Context) {
	var b bindings.UUser
	id := c.Param("id")

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	if ok := bindJSON(c, &b); !ok {
		return
	}

	view, _ := models.UpdateUser(&b, authedUser)
	c.JSON(view.Status, view.View)
}

// DeleteUser is a route handler that is used for user deletion
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	view, _ := models.DeleteUser(authedUser)
	c.Status(view.Status)
}

// VerifyUserEmail is a route handler that is used to verify user's email address using
// a randomly generated cryptographically safe code
func VerifyUserEmail(c *gin.Context) {
	var b bindings.VerifyUserEmail

	id := c.Param("id")

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	if ok := bindJSON(c, &b); !ok {
		return
	}

	view, err := models.VerifyUserEmail(&b, authedUser)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.Status(view.Status)
}

// ReqFriendship is a route handler that is used to request friendship from another
// user in the database
func ReqFriendship(c *gin.Context) {
	var b bindings.Requestee

	id := c.Param("id")

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	if ok := bindJSON(c, &b); !ok {
		return
	}

	view, err := models.ReqFriendship(&b, authedUser)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(view.Status, view.View)
}

// UnReqFriendship is a route handler that is used to delete a friendship request
func UnReqFriendship(c *gin.Context) {
	var b bindings.Requestee

	id := c.Param("id")

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	if ok := bindJSON(c, &b); !ok {
		return
	}

	view, err := models.UnReqFriendship(&b, authedUser)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(view.Status, view.View)
}

// AccFriendship is a route handler that is used to accept a friendship
// request from another user that has been previously requested for friendship
func AccFriendship(c *gin.Context) {
	var b bindings.Requestee

	id := c.Param("id")

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	if ok := bindJSON(c, &b); !ok {
		return
	}

	view, err := models.AccFriendship(&b, authedUser)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(view.Status, view.View)
}

// RejFriendship is a route handler that is used to reject a friendship
// request from another user that has been previously requested for friendship
func RejFriendship(c *gin.Context) {
	var b bindings.Requestee

	id := c.Param("id")

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	if ok := bindJSON(c, &b); !ok {
		return
	}

	view, err := models.RejFriendship(&b, authedUser)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(view.Status, view.View)
}

// CountFriendRequests is a route hander that is used to count user's friend requests
func CountFriendRequests(c *gin.Context) {
	id := c.Param("id")

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	view, _ := models.CountFriendRequests(authedUser)

	c.JSON(view.Status, view.View)
}

// CountFriends is a route hander that is used to count user's friends
func CountFriends(c *gin.Context) {
	id := c.Param("id")

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	view, _ := models.CountFriends(authedUser)

	c.JSON(view.Status, view.View)
}

// ReadFriends is a route hander that is used to get user's friends
func ReadFriends(c *gin.Context) {
	id := c.Param("id")
	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrRequestIsInvalid.Error(),
		})
		return
	}

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	view, _ := models.ReadFriends(page, authedUser)

	c.JSON(view.Status, view.View)
}

// ReadFriendRequests is a route hander that is used to get user's friend
// requests
func ReadFriendRequests(c *gin.Context) {
	id := c.Param("id")
	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrRequestIsInvalid.Error(),
		})
		return
	}

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	view, _ := models.ReadFriendRequests(page, authedUser)

	c.JSON(view.Status, view.View)
}

func areIDAndAuthedUserSame(id string, authedUser string, c *gin.Context) bool {
	if id != authedUser {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": models.ErrUserNotAuthorized.Error(),
		})

		return false
	}

	return true
}
