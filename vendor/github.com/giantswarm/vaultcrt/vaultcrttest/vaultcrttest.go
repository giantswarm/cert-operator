package vaultcrttest

import "github.com/giantswarm/vaultcrt"

type VaultCrtTest struct {
}

func New() *VaultCrtTest {
	return &VaultCrtTest{}
}

func (r *VaultCrtTest) Create(config vaultcrt.CreateConfig) (vaultcrt.CreateResult, error) {
	return vaultcrt.CreateResult{}, nil
}
