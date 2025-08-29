package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neee333ko/errors"
	"github.com/neee333ko/log"
)

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Reference string `json:"reference"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		log.Errorf("%#+v\n", err)
		coder := errors.ParseCoder(err)

		c.JSON(coder.Code(), Response{
			Code:      coder.Code(),
			Message:   coder.Message(),
			Reference: coder.Reference(),
		})
	}

	c.JSON(http.StatusOK, data)
}
