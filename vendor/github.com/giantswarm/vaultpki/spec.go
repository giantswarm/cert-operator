package vaultpki

type Interface interface {
	BackendExists(ID string) (bool, error)
	CAExists(ID string) (bool, error)
	CreateBackend(ID string) error
	CreateCA(ID string) error
	DeleteBackend(ID string) error
}
