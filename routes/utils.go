package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib/models"
)

// ErrRequestIsInvalid is returned when the provided request could not
// be handled due to validation or parsing errors
var ErrRequestIsInvalid = errors.New("Request is invalid")

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
