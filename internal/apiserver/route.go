package apiserver

import (
	"github.com/gin-gonic/gin"
	policyv1 "github.com/neee333ko/IAM/internal/apiserver/controller/v1/policy"
	secretv1 "github.com/neee333ko/IAM/internal/apiserver/controller/v1/secret"
	userv1 "github.com/neee333ko/IAM/internal/apiserver/controller/v1/user"
	"github.com/neee333ko/IAM/internal/apiserver/store/mysql"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
	"github.com/neee333ko/component-base/pkg/core"
	"github.com/neee333ko/errors"
)

func InitRoute(engine *gin.Engine) {
	InitAuth()

	engine.POST("/login", JWTAuth.Login())
	engine.POST("/logout", JWTAuth.Logout())
	engine.GET("/refresh", JWTAuth.Refresh())

	engine.NoRoute(AutoAuth.AuthFunc(), func(ctx *gin.Context) {
		core.WriteResponse(ctx, errors.WithCode(code.ErrPageNotFound, ""), nil)
	})

	store := mysql.GetMysqlInsOr(nil)

	v1 := engine.Group("/v1")

	{
		users := v1.Group("/users")
		{
			UserController := userv1.NewUserController(store)

			users.POST("", UserController.Create)
			users.Use(AutoAuth.AuthFunc(), middleware.Validation())
			users.PUT("/:name", UserController.Update)
			users.PUT("/:name/change-password", UserController.ChangePassword)
			users.GET("/:name", UserController.Get)
			users.GET("", UserController.List)
			users.DELETE("/:name", UserController.Delete)
			users.DELETE("", UserController.DeleteCollection)
		}

		secrets := v1.Group("/secrets")
		{
			SecretController := secretv1.NewSecretController(store)

			secrets.Use(AutoAuth.AuthFunc(), middleware.Validation())

			secrets.POST("", SecretController.Create)
			secrets.GET("/:username", SecretController.Get)
			secrets.DELETE("/name", SecretController.Delete)
			secrets.DELETE("", SecretController.DeleteCollection)
			secrets.GET("", SecretController.List)
		}

		policy := v1.Group("/policy")
		{
			policyController := policyv1.NewPolicyController(store)

			policy.Use(AutoAuth.AuthFunc(), middleware.Validation())

			policy.POST("", policyController.Create)
			policy.PUT("/:name", policyController.Update)
			policy.GET(":username", policyController.Get)
			policy.GET("", policyController.List)
			policy.DELETE("/:name", policyController.Delete)
			policy.DELETE("", policyController.DeleteCollection)
		}
	}
}
