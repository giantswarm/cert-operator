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

func (r *VaultRoleTest) Search(config vaultrole.SearchConfig) (vaultrole.Role, error) {
	return vaultrole.Role{}, nil
}

func (r *VaultRoleTest) Update(config vaultrole.UpdateConfig) error {
	return nil
}
