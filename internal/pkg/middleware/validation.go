package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/apiserver/store"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

var KeyAdmin string = "admin"

func Validation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isadmin, err := isAdmin(ctx)
		username := ctx.GetString(KeyUsername)
		if err != nil {
			core.WriteResponse(ctx, errors.WithCode(code.ErrUnknown, ""), nil)
			ctx.Abort()
			return
		}

		if !isadmin {
			switch ctx.FullPath() {
			case "/v1/users/:name", "/v1/users/:name/change-password":
				name := ctx.Param("name")
				if name != username {
					core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, ""), nil)
					ctx.Abort()
					return
				}
			case "/v1/users", "/v1/secrets", "/v1/policy":
				if ctx.Request.Method == "GET" || ctx.Request.Method == "DELETE" {
					core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, ""), nil)
					ctx.Abort()
					return
				}
			case "/v1/secrets/:username", "/v1/policy/:username":
				un := ctx.Param("username")
				if username != un {
					core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, ""), nil)
					ctx.Abort()
					return
				}
			}
		}
	}
}

func isAdmin(ctx *gin.Context) (bool, error) {
	username := ctx.GetString(KeyUsername)
	user, err := store.Client().NewUserStore().Get(ctx, username, &metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	if user.IsAdmin == 1 {
		ctx.Set(KeyAdmin, 1)
		return true, nil
	}

	return false, nil
}
