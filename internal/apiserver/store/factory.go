package store

var client Factory

type Factory interface {
	NewUserStore() UserStore
	NewSecretStore() SecretStore
	NewPolicyStore() PolicyStore
}

func SetClient(f Factory) {
	client = f
}

func Client() Factory {
	return client
}
