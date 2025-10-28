package policy

import (
	servv1 "github.com/neee333ko/IAM/internal/apiserver/service/v1"
	"github.com/neee333ko/IAM/internal/apiserver/store"
)

type PolicyController struct {
	service servv1.Service
}

func NewPolicyController(s store.Factory) *PolicyController {
	return &PolicyController{service: servv1.NewService(s)}
}
