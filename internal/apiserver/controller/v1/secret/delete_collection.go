package secret

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
)

type secretNames struct {
	Names []string `json:"names"`
}

func (sc *SecretController) DeleteCollection(ctx *gin.Context) {
	log.L(ctx).Info("secret delete collection function called...")

	sn := new(secretNames)

	if err := ctx.BindJSON(sn); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, ""), nil)
		return
	}

	if err := sc.service.SecretServ().DeleteCollection(ctx, sn.Names, &metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
