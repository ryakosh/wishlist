package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib/bindings"
	"github.com/ryakosh/wishlist/lib/models"
)

// TODO: DRY body binding and error handling

func main() {
	r := router()

	r.Run(":8080")
}

func router() *gin.Engine {
	r := gin.Default()

	r.POST("/login", loginUser)
	r.POST("/users", createUser)
	r.GET("/users", readUser)
	r.PUT("/users", models.Authenticate(), updateUser)
	r.DELETE("/users", models.Authenticate(), deleteUser)

	return r
}

// TODO: Don't forget about CSRF attacks
func loginUser(c *gin.Context) {
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

func createUser(c *gin.Context) {
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

func readUser(c *gin.Context) {
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

func updateUser(c *gin.Context) {
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

func deleteUser(c *gin.Context) {
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
