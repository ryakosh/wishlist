package lib

import (
	"context"

	"github.com/gin-gonic/gin"
)

type key int

const (
	ginCtxKey key = iota
)

func GinCtxFromCtx(ctx context.Context) *gin.Context {
	ginCtx := ctx.Value(ginCtxKey)
	if ginCtx == nil {
		LogError(LFatal, "Could not retrieve gin.Context", nil)
	}

	c, ok := ginCtx.(*gin.Context)
	if !ok {
		LogError(LFatal, "Gin.Context has wrong type", nil)
	}
	return c
}

func GinCtxToCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginCtxKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
