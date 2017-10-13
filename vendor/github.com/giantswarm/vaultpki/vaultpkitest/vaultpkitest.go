package vaultpkitest

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
