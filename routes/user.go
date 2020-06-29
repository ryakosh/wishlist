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
	c.JSON(http.StatusOK, gin.H{
		"token": token,
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
	c.JSON(http.StatusCreated, view)
}

// ReadUser is a route hander that is used to get general information about a user
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
	c.JSON(http.StatusOK, view)
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

	view := models.UpdateUser(&b, authedUser)
	c.JSON(http.StatusOK, view)
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

	models.DeleteUser(authedUser)
	c.Status(http.StatusOK)
}

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

	err := models.VerifyUserEmail(&b, authedUser)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.Status(http.StatusOK)
}

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

	c.JSON(http.StatusOK, view)
}

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

	c.JSON(http.StatusOK, view)
}

func CountFriends(c *gin.Context) {
	id := c.Param("id")

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if !areIDAndAuthedUserSame(id, authedUser, c) {
		return
	}

	view := models.CountFriends(authedUser)

	c.JSON(http.StatusOK, view)
}

func ReadFriends(c *gin.Context) {
	id := c.Param("id")
	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 64)
	if err != nil {
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

	view := models.ReadFriends(page, authedUser)

	c.JSON(http.StatusOK, view)
}

func ReadFriendRequests(c *gin.Context) {
	id := c.Param("id")
	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 64)
	if err != nil {
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

	view := models.ReadFriendRequests(page, authedUser)

	c.JSON(http.StatusOK, view)
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
