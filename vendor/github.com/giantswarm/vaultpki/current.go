package vaultpki

import (
	"fmt"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/vaultpki/key"
)

func (p *VaultPKI) BackendExists(ID string) (bool, error) {
	fmt.Printf("start vaultpki\n")
	defer fmt.Printf("end vaultpki\n")

	fmt.Printf("1\n")
	// Check if a PKI for the given cluster ID exists.
	mounts, err := p.vaultClient.Sys().ListMounts()
	if IsNoVaultHandlerDefined(err) {
		fmt.Printf("2\n")
		return false, nil
	} else if err != nil {
		fmt.Printf("3\n")
		return false, microerror.Mask(err)
	}
	fmt.Printf("key.ListMountsPath(ID): %#v\n", key.ListMountsPath(ID))
	fmt.Printf("mounts: %#v\n", mounts)
	mountOutput, ok := mounts[key.ListMountsPath(ID)+"/"]
	fmt.Printf("mountOutput: %#v\n", mountOutput)
	if !ok || mountOutput.Type != "pki" {
		fmt.Printf("4\n")
		return false, nil
	}

	fmt.Printf("5\n")
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
