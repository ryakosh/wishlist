package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib"
)

var accessLog *log.Logger

// AccessLogger is a middleware that is used to log information
// about client's remote address, request's method and request's path
func AccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		remoteAddr := c.Request.RemoteAddr
		method := c.Request.Method
		path := c.Request.URL.Path

		accessLog.Printf("%s - %s - %s\n", remoteAddr, method, path)
	}
}

func init() {
	accessLogFile, err := os.OpenFile("access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		lib.LogError(lib.LFatal, "Could not create log file", err)

	}

	accessLog = log.New(accessLogFile, "", log.LstdFlags)
}
