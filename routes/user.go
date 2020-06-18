package routes

import (
	"net/http"

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
		c.JSON(http.StatusOK, gin.H{
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
		c.JSON(http.StatusOK, gin.H{
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

	if ok := bindJSON(c, &b); !ok {
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
			"error": models.ErrUserNotAuthorized.Error(),
		})

		return
	}

	models.DeleteUser(authedUser.(string))
	c.Status(http.StatusOK)
}
