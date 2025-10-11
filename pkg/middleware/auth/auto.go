package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/pkg/middleware"
	"github.com/neee333ko/component-base/pkg/core"
	"github.com/neee333ko/errors"
)

type AutoAuth struct {
	basic middleware.Auth
	jwt   middleware.Auth
}

var _ middleware.Auth = &AutoAuth{}

func NewAutoAuth(basic, jwt middleware.Auth) *AutoAuth {
	return &AutoAuth{
		basic: basic,
		jwt:   jwt,
	}
}

func (aa *AutoAuth) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := ctx.GetHeader("Authorization")

		strs := strings.SplitN(h, " ", 2)

		if len(strs) != 2 {
			core.WriteResponse(ctx, errors.WithCode(code.ErrInvalidAuthHeader, "Invalid AuthHeader"), nil)
			ctx.Abort()
			return
		}

		var authOp middleware.AuthOperator

		switch strs[0] {
		case "Basic":
			authOp.SetAuth(aa.basic)
		case "Bearer":
			authOp.SetAuth(aa.jwt)
		default:
			core.WriteResponse(ctx, errors.WithCode(code.ErrInvalidAuthHeader, "Invalid AuthHeader"), nil)
			ctx.Abort()
			return
		}

		authOp.Operate()
	}
}
