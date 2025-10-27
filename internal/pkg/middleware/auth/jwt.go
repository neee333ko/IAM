package auth

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
)

type JwtAuth struct {
	m *jwt.GinJWTMiddleware
}

var _ middleware.Auth = &JwtAuth{}

func NewJwtAuth(m *jwt.GinJWTMiddleware) *JwtAuth {
	return &JwtAuth{
		m: m,
	}
}

func (j *JwtAuth) AuthFunc() gin.HandlerFunc {
	return j.m.MiddlewareFunc()
}
