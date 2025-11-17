package authzserver

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/authzserver/controller/v1/authz"
	"github.com/neee333ko/IAM/internal/authzserver/load/cache"
	"github.com/neee333ko/IAM/internal/authzserver/service/v1/authorize"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/internal/pkg/middleware/auth"
	"github.com/neee333ko/component-base/pkg/core"
	"github.com/neee333ko/errors"
)

func InitRoute(engine *gin.Engine) {
	cache := cache.GetCacheInsOr()

	authz := authz.NewAuthzController(authorize.NewAuthorize(cache))
	cacheAuth := auth.NewCacheAuth()

	engine.NoRoute(cacheAuth.AuthFunc(), func(ctx *gin.Context) {
		core.WriteResponse(ctx, errors.WithCode(code.ErrPageNotFound, ""), nil)
	})

	v1 := engine.Group("/v1")
	{
		v1.Use(cacheAuth.AuthFunc())
		v1.POST("/authz", authz.Authorize)
	}
}
