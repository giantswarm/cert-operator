package vaultpki

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/vaultpki/key"
)

func (p *VaultPKI) BackendExists(ID string) (bool, error) {
	// Check if a PKI for the given cluster ID exists.
	mounts, err := p.vaultClient.Sys().ListMounts()
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}
	mountOutput, ok := mounts[key.ListMountsPath(ID)]
	if !ok || mountOutput.Type != "pki" {
		return false, nil
	}

	return true, nil
}

func (p *VaultPKI) CAExists(ID string) (bool, error) {
	// Check if a root CA for the given cluster ID exists.
	secret, err := p.vaultClient.Logical().Read(key.ReadCAPath(ID))
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}

	// If the secret is nil, the CA has not been generated.
	if secret == nil {
		return false, nil
	}

	dataCertificate, ok := secret.Data["certificate"]
	if ok && dataCertificate == "" {
		return false, nil
	}
	dataError, ok := secret.Data["error"]
	if ok && dataError != "" {
		return false, nil
	}

	return true, nil
}
