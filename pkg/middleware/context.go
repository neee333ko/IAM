package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/log"
)

func Context(ctx *gin.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(log.KeyRequestID, ctx.GetHeader(KeyRequestID))
		ctx.Set(log.KeyUsername, ctx.GetString(KeyUsername))
	}
}
