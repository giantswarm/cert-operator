package vaultpki

import (
	vaultapi "github.com/hashicorp/vault/api"
)

type Interface interface {
	BackendExists(ID string) (bool, error)
	CAExists(ID string) (bool, error)
	CreateBackend(ID string) error
	CreateCA(ID string) error
	DeleteBackend(ID string) error
	GetBackend(ID string) (*vaultapi.MountOutput, error)
	GetCACertificate(ID string) (string, error)
}
