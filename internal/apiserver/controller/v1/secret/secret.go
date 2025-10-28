package secret

import (
	servv1 "github.com/neee333ko/IAM/internal/apiserver/service/v1"
	"github.com/neee333ko/IAM/internal/apiserver/store"
)

type SecretController struct {
	service servv1.Service
}

func NewSecretController(s store.Factory) *SecretController {
	return &SecretController{service: servv1.NewService(s)}
}
