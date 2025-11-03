package user

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/log"
)

func (uc *UserController) Delete(ctx *gin.Context) {
	log.L(ctx).Info("user delete function called...")

	username := ctx.Param("name")

	if err := uc.service.UserServ().Delete(ctx, username, &metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
