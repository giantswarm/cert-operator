package vaultpkitest

import (
	vaultapi "github.com/hashicorp/vault/api"
)

type VaultPKITest struct {
}

func New() *VaultPKITest {
	return &VaultPKITest{}
}

func (p *VaultPKITest) BackendExists(ID string) (bool, error) {
	return false, nil
}

func (p *VaultPKITest) CAExists(ID string) (bool, error) {
	return false, nil
}

func (p *VaultPKITest) CreateBackend(ID string) error {
	return nil
}

func (p *VaultPKITest) CreateCA(ID string) error {
	return nil
}

func (p *VaultPKITest) DeleteBackend(ID string) error {
	return nil
}

func (p *VaultPKITest) GetBackend(ID string) (*vaultapi.MountOutput, error) {
	return nil, nil
}

func (p *VaultPKITest) GetCACertificate(ID string) (string, error) {
	return "", nil
}

func (p *VaultPKITest) ListBackends() ([]*vaultapi.MountOutput, error) {
	return nil, nil
}
