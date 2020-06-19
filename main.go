package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib/models"
	"github.com/ryakosh/wishlist/routes"
)

var r *gin.Engine

func main() {
	r.Run(":8080")
}

func init() {
	r = gin.Default()

	r.POST("/login", routes.LoginUser)
	users := r.Group("/users")
	users.POST("", routes.CreateUser)
	users.GET("/:id", routes.ReadUser)
	users.PUT("", models.Authenticate(), routes.UpdateUser)
	users.DELETE("", models.Authenticate(), routes.DeleteUser)

	wishes := r.Group("/wishes")
	wishes.POST("", models.Authenticate(), routes.CreateWish)
	wishes.GET("/:id", routes.ReadWish)
	wishes.PUT("/:id", models.Authenticate(), routes.UpdateWish)
	wishes.DELETE("/:id", models.Authenticate(), routes.DeleteWish)
}
