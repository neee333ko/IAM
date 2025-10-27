package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/errors"
	"golang.org/x/time/rate"
)

var limitExceedError = errors.New("exceed server limit")

func Limit(eventsPerSecond float64, maxBurstSize int) gin.HandlerFunc {
	l := rate.NewLimiter(rate.Limit(eventsPerSecond), maxBurstSize)

	return func(ctx *gin.Context) {
		if l.Allow() {
			ctx.Next()

			return
		}

		ctx.AbortWithError(429, limitExceedError)
	}
}
