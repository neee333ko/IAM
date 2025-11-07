package auth

import (
	jwt "github.com/appleboy/gin-jwt/v3"
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

func (j *JwtAuth) Login() gin.HandlerFunc {
	return j.m.LoginHandler
}

func (j *JwtAuth) Logout() gin.HandlerFunc {
	return j.m.LogoutHandler
}

func (j *JwtAuth) Refresh() gin.HandlerFunc {
	return j.m.RefreshHandler
}
