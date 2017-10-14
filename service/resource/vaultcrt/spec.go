package vaultcrt

type Secret struct {
	BackendExists bool
	CAExists    bool
	IsPolicyCreated  bool
	IsRoleCreated    bool
}
