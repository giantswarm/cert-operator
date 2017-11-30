package vaultpkiv2

import vaultapi "github.com/hashicorp/vault/api"

type VaultPKIState struct {
	Backend       *vaultapi.MountOutput
	CACertificate string
}
