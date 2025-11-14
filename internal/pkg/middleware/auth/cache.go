package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/neee333ko/IAM/internal/authzserver/load/cache"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
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

		_, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			claims := t.Claims.(jwt.MapClaims)
			secretID := claims["id"].(string)

			key, err := c.cache.GetSecret(secretID)
			if err != nil {
				return nil, err
			}

			ctx.Set(middleware.KeyUsername, key.Username)

			return key.SecretKey, nil
		}, nil)
		if err != nil {
			core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.Error()), nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
