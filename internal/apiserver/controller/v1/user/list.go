package user

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	v1 "github.com/neee333ko/api/apiserver/v1"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

func (uc *UserController) List(ctx *gin.Context) {
	var options *metav1.ListOptions = new(metav1.ListOptions)

	if err := ctx.BindQuery(options); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	var list *v1.UserList = new(v1.UserList)
	list, err := uc.service.UserServ().List(options)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, errors.WithCode(code.ErrSuccess, ""), list)
}
