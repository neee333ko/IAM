package user

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	v1 "github.com/neee333ko/api/apiserver/v1"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

func (uc *UserController) Update(ctx *gin.Context) {
	var user *v1.User = new(v1.User)

	if err := ctx.BindJSON(user); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	username := ctx.Param("name")
	u, err := uc.service.UserServ().Get(ctx, username, &metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	u.Status = user.Status
	u.Nickname = user.Nickname
	u.Email = user.Email
	u.Phone = user.Phone

	if err := u.ValidateUpdate(); len(err) != 0 {
		core.WriteResponse(ctx, errors.WithCode(code.ErrValidation, err.ToAggregate().Error()), nil)
		return
	}

	if err := uc.service.UserServ().Update(ctx, u, &metav1.UpdateOptions{}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, u)
}
