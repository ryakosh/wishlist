package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib"
	"github.com/ryakosh/wishlist/lib/models"
	"github.com/ryakosh/wishlist/routes"
)

var r *gin.Engine

func main() {
	r.Run(":8080")
}

func init() {
	serverLog, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		lib.LogError(lib.LFatal, "Could not create log file", err)

	}
	log.SetOutput(serverLog)

	gin.DisableConsoleColor()
	ginLog, err := os.OpenFile("gin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		lib.LogError(lib.LFatal, "Could not create log file", err)
	}
	gin.DefaultWriter = ginLog

	r = gin.Default()
	r.Use(AccessLogger())

	r.POST("/login", routes.LoginUser)
	users := r.Group("/users")
	users.POST("", routes.CreateUser)
	users.GET("/:id", routes.ReadUser)
	users.PUT("/:id", models.Authenticate(), routes.UpdateUser)
	users.DELETE("/:id", models.Authenticate(), routes.DeleteUser)
	users.PUT("/:id/verify_email", models.Authenticate(), routes.VerifyUserEmail)

	friendRequests := users.Group("/:id/friend_requests")
	friendRequests.GET("", models.Authenticate(), models.RequireEmailVerification(), routes.ReadFriendRequests)
	friendRequests.GET("/count", models.Authenticate(), models.RequireEmailVerification(), routes.CountFriendRequests)

	friends := users.Group("/:id/friends")
	friends.GET("", models.Authenticate(), models.RequireEmailVerification(), routes.ReadFriends)
	friends.GET("/count", models.Authenticate(), models.RequireEmailVerification(), routes.CountFriends)
	friends.PUT("/send_request", models.Authenticate(), models.RequireEmailVerification(), routes.ReqFriendship)
	friends.DELETE("/undo_request", models.Authenticate(), models.RequireEmailVerification(), routes.UnReqFriendship)
	friends.PUT("/accept_request", models.Authenticate(), models.RequireEmailVerification(), routes.AccFriendship)
	friends.DELETE("/reject_request", models.Authenticate(), models.RequireEmailVerification(), routes.RejFriendship)

	wishes := r.Group("/wishes")
	wishes.POST("", models.Authenticate(), models.RequireEmailVerification(), routes.CreateWish)
	wishes.GET("/:id", routes.ReadWish)
	wishes.PUT("/:id", models.Authenticate(), routes.UpdateWish)
	wishes.DELETE("/:id", models.Authenticate(), routes.DeleteWish)
	wishes.PUT("/:id/add_fulfiller", models.Authenticate(), models.RequireEmailVerification(), routes.AddWantToFulfill)
	wishes.PUT("/:id/add_claimer", models.Authenticate(), routes.AddClaimer)
	wishes.PUT("/:id/accept_claimer", models.Authenticate(), routes.AddWantToFulfill)
	wishes.PUT("/:id/reject_claimer", models.Authenticate(), routes.RejectClaimer)
	wishes.PUT("/:id/accept_claimer", models.Authenticate(), routes.AddWantToFulfill)
	wishes.GET("/:id/read_fulfillers", models.Authenticate(), routes.ReadFulfillers)
	wishes.GET("/:id/read_claimers", models.Authenticate(), routes.ReadClaimers)
	wishes.GET("/:id/count_fulfillers", routes.CountWantToFulfill)
}
