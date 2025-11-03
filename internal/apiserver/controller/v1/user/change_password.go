package user

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/component-base/pkg/auth"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
)

type password struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (uc *UserController) ChangePassword(ctx *gin.Context) {
	log.L(ctx).Info("user change password function called...")

	var pw *password = new(password)

	if err := ctx.ShouldBindJSON(pw); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	username := ctx.Param("name")

	u, err := uc.service.UserServ().Get(ctx, username, &metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	if err := u.Compare(pw.OldPassword); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrPasswordIncorrect, ""), nil)
		return
	}

	u.Password = pw.NewPassword

	if err := u.ValidateUpdatePassword(); len(err) != 0 {
		core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.ToAggregate().Error()), nil)
		return
	}

	u.Password, _ = auth.Encrypt(u.Password)

	if err := uc.service.UserServ().Update(ctx, u, &metav1.UpdateOptions{}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, u)
}
