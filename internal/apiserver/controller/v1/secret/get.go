package secret

import (
	"github.com/gin-gonic/gin"
	"github.com/neee333ko/IAM/internal/pkg/code"
	"github.com/neee333ko/IAM/internal/pkg/middleware"
	"github.com/neee333ko/component-base/pkg/core"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/errors"
)

func (sc *SecretController) Get(ctx *gin.Context) {
	username := ctx.Param(":username")

	sl, err := sc.service.SecretServ().Get(ctx, username, &metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	if ctx.GetInt(middleware.KeyAdmin) != 1 {
		for i := range sl.Items {
			sl.Items[i].SecretKey = ""
		}
	}

	core.WriteResponse(ctx, errors.WithCode(code.ErrSuccess, ""), sl)
}
