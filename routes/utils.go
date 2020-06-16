package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrBodyIsInvalid is returned when the provided request body could not
// be handled due to validation or parsing errors
var ErrBodyIsInvalid = errors.New("Request body is invalid")

func bindJSON(c *gin.Context, bindTo interface{}) bool {
	if err := c.ShouldBindJSON(bindTo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrBodyIsInvalid.Error(),
		})

		return false
	}

	return true
}
