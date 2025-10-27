package secret

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
	v1 "github.com/neee333ko/api/apiserver/v1"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/component-base/pkg/util/idutil"
	"github.com/neee333ko/errors"
)

var maxSecretCount = 10

func (sc *SecretController) Create(ctx *gin.Context) {
	var secret *v1.Secret = new(v1.Secret)

	if err := ctx.BindJSON(secret); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	username := ctx.GetString(middleware.KeyUsername)

	if username != secret.Username && ctx.GetInt(middleware.KeyAdmin) != 1 {
		core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, ""), nil)
		return
	}

	sl, err := sc.service.SecretServ().Get(ctx, secret.Username, &metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	if len(sl.Items) >= maxSecretCount {
		core.WriteResponse(ctx, errors.WithCode(code.ErrReachMaxCount, ""), nil)
		return
	}

	secret.SecretID = idutil.NewSecretID()
	secret.SecretKey = idutil.NewSecretKey()

	if err := secret.Validate(); len(err) != 0 {
		core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.ToAggregate().Error()), nil)
		return
	}

	if err := sc.service.SecretServ().Create(ctx, secret, &metav1.CreateOptions{}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, errors.WithCode(code.ErrSuccess, ""), secret)
}
