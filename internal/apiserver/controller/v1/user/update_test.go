package user

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	servv1 "github.com/neee333ko/IAM/internal/apiserver/service/v1"
	v1 "github.com/neee333ko/api/apiserver/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"go.uber.org/mock/gomock"
)

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockserv := servv1.NewMockService(ctrl)
	mockuserserv := servv1.NewMockUserServ(ctrl)

	mockserv.EXPECT().UserServ().Return(mockuserserv).Times(2)
	mockuserserv.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	user := &v1.User{ObjectMeta: metav1.ObjectMeta{Name: "neee333ko"}, Nickname: "neee333ko", Password: "Wto1260644864!", Email: "1260644864@qq.com"}
	mockuserserv.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, make(field.ErrorList, 0).ToAggregate())

	buffer := bytes.NewBufferString(`{"nickname":"neeeko","email":"wto1260644864@qq.com"}`)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	request := httptest.NewRequest("PUT", "/v1/users/:name", buffer)
	request.Header.Set("Content-Type", "application/json")
	c.Request = request
	c.Params = append(c.Params, gin.Param{Key: "name", Value: "neee333ko"})

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

			uc.Update(tt.args.ctx)
		})
	}
}
