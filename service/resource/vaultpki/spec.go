package vaultpki

type VaultPKIState struct {
	BackendExists bool
	CAExists    bool
	IsPolicyCreated  bool
	IsRoleCreated    bool
}
