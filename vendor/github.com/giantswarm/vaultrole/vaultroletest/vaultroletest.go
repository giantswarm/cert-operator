package vaultroletest

import "github.com/giantswarm/vaultrole"

type VaultRoleTest struct {
}

func New() *VaultRoleTest {
	return &VaultRoleTest{}
}

func (r *VaultRoleTest) Create(config vaultrole.CreateConfig) error {
	return nil
}

func (r *VaultRoleTest) Exists(config vaultrole.ExistsConfig) (bool, error) {
	return false, nil
}
