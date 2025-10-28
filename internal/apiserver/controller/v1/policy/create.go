package policy

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
	v1 "github.com/neee333ko/api/apiserver/v1"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

func (pc *PolicyController) Create(ctx *gin.Context) {
	username := ctx.GetString(middleware.KeyUsername)

	var policy *v1.Policy = new(v1.Policy)

	if err := ctx.BindJSON(policy); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	if policy.Username != username && ctx.GetInt(middleware.KeyAdmin) != 1 {
		core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, ""), nil)
		return
	}

	if err := policy.Validate(); len(err) != 0 {
		core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.ToAggregate().Error()), nil)
		return
	}

	if err := pc.service.PolicyServ().Create(ctx, policy, &metav1.CreateOptions{}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
