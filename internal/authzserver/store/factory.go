package store

var client Factory

type Factory interface {
	NewSecretStore() SecretStore
	NewPolicyStore() PolicyStore
}

func SetClient(f Factory) {
	client = f
}

func Client() Factory {
	return client
}
