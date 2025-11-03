package user

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
)

type userNames struct {
	Names []string `json:"names"`
}

func (uc *UserController) DeleteCollection(ctx *gin.Context) {
	log.L(ctx).Info("user delete collection function called...")

	var users *userNames = new(userNames)

	if err := ctx.ShouldBindJSON(users); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	if err := uc.service.UserServ().DeleteCollection(ctx, users.Names, &metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
