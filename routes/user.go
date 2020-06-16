package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib/bindings"
	"github.com/ryakosh/wishlist/lib/models"
)

// LoginUser is a route handler that is used for user authentication,
// it does so by setting a cookie that contains a token
func LoginUser(c *gin.Context) { // TODO: Don't forget about CSRF attacks
	var b bindings.LoginUser

	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not bind the provided json",
		})

		return
	}

	token, err := models.LoginUser(&b)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})

		return
	}

	// TODO: User should not be able to authenticate when they are authenticated already
	// TODO: Update domain when deploying
	// TODO: Set secure cookie to true in production
	c.SetCookie(models.TokenCookieKey, token, -1, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}

// CreateUser is a route handler that is used to create/register a new user
func CreateUser(c *gin.Context) {
	var b bindings.CUser

	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not bind the provided json",
		})

		return
	}

	view, err := models.CreateUser(&b)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})

		return
	}
	c.JSON(http.StatusCreated, view)
}

// ReadUser is a route hander that is used to get general information about a user
func ReadUser(c *gin.Context) {
	var b bindings.RUser

	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not bind the provided json",
		})

		return
	}

	view, err := models.ReadUser(&b)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, view)
}

// UpdateUser is a route handler that is used to update general information about a user
func UpdateUser(c *gin.Context) {
	var b bindings.UUser

	authedUser, ok := c.Get(models.UserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": models.ErrUserNotAuthorized.Error(),
		})

		return
	}

	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not bind the provided json",
			"err":   err.Error(),
		})

		return
	}

	view := models.UpdateUser(&b, authedUser.(string))
	c.JSON(http.StatusOK, view)
}

// DeleteUser is a route handler that is used for user deletion
func DeleteUser(c *gin.Context) {
	authedUser, ok := c.Get(models.UserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": models.ErrUserNotAuthorized,
		})

		return
	}

	models.DeleteUser(authedUser.(string))
	c.Status(http.StatusOK)
}
