package auth

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/neee333ko/IAM/internal/authzserver/load/cache"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
	pb "github.com/neee333ko/api/proto/v1"
	"github.com/neee333ko/component-base/pkg/core"
	"github.com/neee333ko/errors"
)

type CacheAuth struct {
	cache *cache.Cache
}

func NewCacheAuth() *CacheAuth {
	return &CacheAuth{cache: cache.GetCacheInsOr()}
}

func (c *CacheAuth) AuthFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authString := ctx.GetHeader("Authorization")
		parts := strings.Split(authString, " ")
		if len(parts) != 2 {
			core.WriteResponse(ctx, errors.WithCode(code.ErrInvalidAuthHeader, ""), nil)
			return
		}

		token := parts[1]

		var secretInfo *pb.SecretInfo = new(pb.SecretInfo)
		t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			secretID := t.Header["kid"].(string)

			var err error
			secretInfo, err = c.cache.GetSecret(secretID)
			if err != nil {
				return nil, err
			}

			return []byte(secretInfo.SecretKey), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !t.Valid {
			core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), nil)
			ctx.Abort()
			return
		}

		if expire(secretInfo.Expires) {
			core.WriteResponse(ctx, errors.WithCode(code.ErrExpired, ""), nil)
			ctx.Abort()
			return
		}

		ctx.Set(middleware.KeyUsername, secretInfo.Username)

		ctx.Next()
	}
}

func expire(e int64) bool {
	if e >= 1 {
		return time.Now().After(time.Unix(e, 0))
	}

	return true
}
