package user

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
)

func (uc *UserController) Get(ctx *gin.Context) {
	username := ctx.Param("name")

	user, err := uc.service.UserServ().Get(ctx, username, &metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, user)
}
