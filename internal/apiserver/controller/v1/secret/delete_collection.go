package secret

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	v1 "github.com/neee333ko/api/apiserver/v1"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

func (sc *SecretController) DeleteCollection(ctx *gin.Context) {
	var secretlist *v1.SecretList = new(v1.SecretList)

	if err := ctx.BindJSON(secretlist); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, ""), nil)
		return
	}

	if err := sc.service.SecretServ().DeleteCollection(secretlist, &metav1.DeleteOptions{}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, errors.WithCode(code.ErrSuccess, ""), nil)
}
