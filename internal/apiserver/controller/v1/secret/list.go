package secret

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
)

func (sc *SecretController) List(ctx *gin.Context) {
	log.L(ctx).Info("secret list function called...")

	var options *metav1.ListOptions = new(metav1.ListOptions)

	if err := ctx.BindQuery(options); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	list, err := sc.service.SecretServ().List(ctx, options)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, list)
}
