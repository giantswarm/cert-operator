package vaultcrttest

import "github.com/giantswarm/vaultcrt"

type VaultCrtTest struct {
}

func New() *VaultCrtTest {
	return &VaultCrtTest{}
}

func (r *VaultCrtTest) Create(config vaultcrt.CreateConfig) (vaultcrt.CreateResult, error) {
	result := vaultcrt.CreateResult{
		CA:           "test CA",
		Crt:          "test crt",
		Key:          "test key",
		SerialNumber: "test serial number",
	}
	return result, nil
}
