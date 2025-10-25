package user

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	v1 "github.com/neee333ko/api/apiserver/v1"
	"github.com/neee333ko/component-base/pkg/auth"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

func (uc *UserController) Create(ctx *gin.Context) {
	var user *v1.User = new(v1.User)
	if err := ctx.BindJSON(user); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	if err := user.ValidateCreate(); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.ToAggregate().Error()), nil)
		return
	}

	user.Password, _ = auth.Encrypt(user.Password)
	user.Status = 1

	if err := uc.service.UserServ().Create(user, &metav1.CreateOptions{}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, errors.WithCode(code.ErrSuccess, ""), nil)
}
