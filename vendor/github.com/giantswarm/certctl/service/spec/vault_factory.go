package spec

import (
	vault "github.com/hashicorp/vault/api"
)

// VaultFactory implements a factory that is able to create Vault clients.
type VaultFactory interface {
	// NewClient creates a new Vault client configured with an admin token.
	NewClient() (*vault.Client, error)
}
