package authz

import (
	"github.com/gin-gonic/gin"
	servv1 "github.com/neee333ko/IAM/internal/authzserver/service/v1"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
	v1 "github.com/neee333ko/api/authzserver/v1"
	"github.com/neee333ko/component-base/pkg/core"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
	"github.com/ory/ladon"
)

type AuthzController struct {
	authorizer servv1.Authorizer
}

func NewAuthzController(authorizer servv1.Authorizer) *AuthzController {
	return &AuthzController{
		authorizer: authorizer,
	}
}

func (c *AuthzController) Authorize(ctx *gin.Context) {
	log.Info("authz fn called...")

	var req *ladon.Request = new(ladon.Request)
	if err := ctx.ShouldBindJSON(req); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}
	if req.Context == nil {
		req.Context = ladon.Context{}
	}

	req.Context["username"] = ctx.GetString(middleware.KeyUsername)

	if err := servv1.NewAuthzService(c.authorizer).IsAllowed(ctx, req); err != nil {
		core.WriteResponse(ctx, nil, &v1.Response{Result: "false", Error: err.Error()})
	}

	core.WriteResponse(ctx, nil, &v1.Response{Result: "true", Error: "nil"})
}
