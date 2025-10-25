package v1

type Service interface {
	UserServ() UserServ
	SecretServ() SecretServ
	PolicyServ() PolicyServ
}
