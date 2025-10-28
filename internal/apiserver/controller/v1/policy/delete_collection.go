package policy

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

type policyNames struct {
	Names []string `json:"names"`
}

func (pc *PolicyController) DeleteCollection(ctx *gin.Context) {
	pn := new(policyNames)

	if err := ctx.BindJSON(pn); err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrBind, ""), nil)
		return
	}

	if err := pc.service.SecretServ().DeleteCollection(ctx, pn.Names, &metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
