package user

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	servv1 "github.com/neee333ko/IAM/internal/apiserver/service/v1"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockserv := servv1.NewMockService(ctrl)
	mockuserserv := servv1.NewMockUserServ(ctrl)

	mockserv.EXPECT().UserServ().Return(mockuserserv)
	mockuserserv.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	buffer := bytes.NewBufferString(`{"metadata":{"name":"neee333ko"},"nickname":"neko","password":"Wto1260644864!","email":"1260644864@qq.com"}`)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	request := httptest.NewRequest("POST", "/v1/users", buffer)
	request.Header.Set("Content-Type", "application/json")
	c.Request = request

	type field struct {
		service servv1.Service
	}

	type arg struct {
		ctx *gin.Context
	}

	tests := []struct {
		name   string
		fields field
		args   arg
	}{
		{name: "default", fields: field{service: mockserv}, args: arg{ctx: c}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &UserController{
				service: mockserv,
			}

			uc.Create(tt.args.ctx)
		})
	}
}
