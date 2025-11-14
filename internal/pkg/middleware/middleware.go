package middleware

import "github.com/gin-gonic/gin"

var KeyUsername string = "username"

var Middlewares map[string]gin.HandlerFunc = DefaultMiddlewares()

func DefaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery": gin.Recovery(),
		"logger":   Logger(),
		"limit":    Limit(10, 5),
		"request":  RequestID(),
		"context":  Context(),
		"cors":     Cors(),
		"publish":  Publish(),
	}
}
