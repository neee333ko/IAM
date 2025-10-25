package user

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	v1 "github.com/neee333ko/api/apiserver/v1"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

func (uc *UserController) ChangePassword(ctx *gin.Context) {
	var user *v1.User = new(v1.User)

	if err := ctx.BindJSON(user); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	username := ctx.Param("name")
	user.Name = username

	if err := uc.service.UserServ().Update(user, &metav1.UpdateOptions{}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, errors.WithCode(code.ErrSuccess, ""), nil)
}
