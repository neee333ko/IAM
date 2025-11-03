package policy

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
)

func (pc *PolicyController) List(ctx *gin.Context) {
	log.L(ctx).Info("policy list function called...")

	var options *metav1.ListOptions = new(metav1.ListOptions)

	if err := ctx.ShouldBindQuery(options); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	list, err := pc.service.PolicyServ().List(ctx, options)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, list)
}
