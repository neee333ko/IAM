package secret

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	v1 "github.com/neee333ko/api/apiserver/v1"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

func (sc *SecretController) Create(ctx *gin.Context) {
	var secret *v1.Secret = new(v1.Secret)

	if err := ctx.BindJSON(secret); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	if err := sc.service.SecretServ().Create(secret, &metav1.CreateOptions{}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, errors.WithCode(code.ErrSuccess, ""), nil)
}
