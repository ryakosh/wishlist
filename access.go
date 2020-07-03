package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ryakosh/wishlist/lib"
)

var accessLog *log.Logger

// AccessLogger is a middleware that is used to log information
// about client's remote address, request's method, request's path and
// request's body
func AccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		remoteAddr := c.Request.RemoteAddr
		method := c.Request.Method
		path := c.Request.URL.Path
		body, err := ioutil.ReadAll(c.Request.Body)

		if err != nil {
			lib.LogError(lib.LError, "Could not log request's body to access log", err)
		}

		accessLog.Printf("%s - %s - %s\n%s\n\n", remoteAddr, method, path, string(body))
	}
}

func init() {
	accessLogFile, err := os.OpenFile("access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		lib.LogError(lib.LFatal, "Could not create log file", err)

	}

	accessLog = log.New(accessLogFile, "", log.LstdFlags)
}
