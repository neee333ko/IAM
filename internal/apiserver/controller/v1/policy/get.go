package policy

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/log"
)

func (pc *PolicyController) Get(ctx *gin.Context) {
	log.L(ctx).Info("policy get function called...")

	username := ctx.Param("username")

	policy, err := pc.service.PolicyServ().Get(ctx, username, &metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, policy)
}
