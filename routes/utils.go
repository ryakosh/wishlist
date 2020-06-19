package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
