package secret

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
)

func (sc *SecretController) Delete(ctx *gin.Context) {
	log.L(ctx).Info("secret delete function called...")

	name := ctx.Param("name")
	username := ctx.GetString(middleware.KeyUsername)

	s, err := sc.service.SecretServ().GetSingle(ctx, name, &metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	if s.Username != username && ctx.GetInt(middleware.KeyAdmin) != 1 {
		core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, ""), nil)
		return
	}

	if err := sc.service.SecretServ().Delete(ctx, name, &metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
