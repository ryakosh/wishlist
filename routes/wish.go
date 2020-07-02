package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib/bindings"
	"github.com/ryakosh/wishlist/lib/models"
)

// CreateWish is a route handler that is used to create a new wish
func CreateWish(c *gin.Context) {
	var b bindings.CWish

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if ok := bindJSON(c, &b); !ok {
		return
	}

	view, _ := models.CreateWish(&b, authedUser)
	c.JSON(view.Status, view.View)
}

// ReadWish is a route hander that is used to get general information about a wish
func ReadWish(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrRequestIsInvalid.Error(),
		})
	}

	view, err := models.ReadWish(id)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}
	c.JSON(view.Status, view.View)
}

// UpdateWish is a route handler that is used to update general information about a wish
func UpdateWish(c *gin.Context) {
	var b bindings.UWish

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrRequestIsInvalid.Error(),
		})
	}

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	if ok := bindJSON(c, &b); !ok {
		return
	}

	view, err := models.UpdateWish(id, &b, authedUser)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(view.Status, view.View)
}

// DeleteWish is a route handler that is used for wish deletion
func DeleteWish(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrRequestIsInvalid.Error(),
		})
	}

	authedUser := authedUser(c)
	if authedUser == "" {
		return
	}

	view, err := models.DeleteWish(id, authedUser)
	if err != nil {
		err := err.(*models.RequestError)
		c.JSON(err.Status, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.Status(view.Status)
}
