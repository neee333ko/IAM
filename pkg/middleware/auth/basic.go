package auth

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/pkg/middleware"
	"github.com/neee333ko/component-base/pkg/core"
	"github.com/neee333ko/errors"
)

type BasicAuth struct {
	Compare func(name string, password string) bool
}

var _ middleware.Auth = &BasicAuth{}

func NewBasicAuth(compare func(name string, password string) bool) *BasicAuth {
	return &BasicAuth{
		Compare: compare,
	}
}

func (b *BasicAuth) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		head := ctx.GetHeader("Authorization")
		authorization := strings.SplitN(head, " ", 2)

		if authorization[0] != "Basic" || len(authorization) != 2 {
			core.WriteResponse(ctx, errors.WithCode(code.ErrInvalidAuthHeader, "Invalid AuthHeader"), nil)
			ctx.Abort()
			return
		}

		token, _ := base64.StdEncoding.DecodeString(authorization[1])

		strs := strings.SplitN(string(token), ":", 2)

		if len(strs) != 2 {
			core.WriteResponse(ctx, errors.WithCode(code.ErrInvalidAuthHeader, "Invalid AuthHeader"), nil)
			ctx.Abort()
			return
		}

		if !b.Compare(strs[0], strs[1]) {
			core.WriteResponse(ctx, errors.WithCode(code.ErrPasswordIncorrect, "Password Incorrect"), nil)
			ctx.Abort()
			return
		}

		ctx.Set(middleware.KeyUsername, strs[0])

		ctx.Next()

	}
}
