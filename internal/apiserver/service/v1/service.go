package v1

import "github.com/neee333ko/IAM/internal/apiserver/store"

type Service interface {
	UserServ() UserServ
	SecretServ() SecretServ
	PolicyServ() PolicyServ
}

type SService struct {
	store store.Factory
}

func (s *SService) UserServ() UserServ {
	return &UserService{store: s.store}
}

func (s *SService) SecretServ() SecretServ {
	return &SecretService{store: s.store}
}

func (s *SService) PolicyServ() PolicyServ {
	return &PolicyService{store: s.store}
}
