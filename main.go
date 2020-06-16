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
	r.POST("/users", routes.CreateUser)
	r.GET("/users", routes.ReadUser)
	r.PUT("/users", models.Authenticate(), routes.UpdateUser)
	r.DELETE("/users", models.Authenticate(), routes.DeleteUser)
}
