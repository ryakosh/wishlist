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

	view := models.CreateWish(&b, authedUser.(string))
	c.JSON(http.StatusOK, view)
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
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, view)
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

	view, err := models.UpdateWish(id, &b, authedUser.(string))
	if err == models.ErrUserNotAuthorized {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	} else if err == models.ErrWishNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, view)
}

// DeleteWish is a route handler that is used for wish deletion
func DeleteWish(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrRequestIsInvalid.Error(),
		})
	}

	authedUser, ok := c.Get(models.UserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": models.ErrUserNotAuthorized.Error(),
		})

		return
	}

	err = models.DeleteWish(id, authedUser.(string))
	if err == models.ErrUserNotAuthorized {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	} else if err == models.ErrWishNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.Status(http.StatusOK)
}
