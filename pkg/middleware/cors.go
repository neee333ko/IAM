package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowOriginFunc: func(origin string) bool {
				return origin == "https//github.com"
			},
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders: []string{"Content-Type", "Authorization", "Origin", "Accept"},
			MaxAge:       time.Hour,
		},
	)
}
