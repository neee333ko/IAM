package user

import (
	servv1 "github.com/neee333ko/IAM/internal/apiserver/service/v1"
)

type UserController struct {
	service servv1.Service
}

// func NewUserController(s store.Factory) *UserController {
// 	return &UserController{
// 		service: v1.new,
// 	}
// }
