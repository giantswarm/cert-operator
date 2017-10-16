package vaultpki

import (
	"fmt"

	"github.com/giantswarm/microerror"
	vaultapi "github.com/hashicorp/vault/api"

	"github.com/giantswarm/vaultpki/key"
)

func (p *VaultPKI) BackendExists(ID string) (bool, error) {
	_, err := p.GetBackend(ID)
	if IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}

	return true, nil
}

func (p *VaultPKI) CAExists(ID string) (bool, error) {
	_, err := p.GetCACertificate(ID)
	if IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}

	return true, nil
}

func (p *VaultPKI) GetBackend(ID string) (*vaultapi.MountOutput, error) {
	fmt.Printf("start VaultPKI.GetBackend\n")
	defer fmt.Printf("end VaultPKI.GetBackend\n")

	fmt.Printf("1\n")

	mounts, err := p.vaultClient.Sys().ListMounts()
	if IsNoVaultHandlerDefined(err) {
		fmt.Printf("2\n")
		return nil, microerror.Maskf(notFoundError, "PKI backend for ID '%s'", ID)
	} else if err != nil {
		fmt.Printf("3\n")
		return nil, microerror.Mask(err)
	}

	fmt.Printf("%#v\n", key.ListMountsPath(ID))
	fmt.Printf("%#v\n", mounts)

	mountOutput, ok := mounts[key.ListMountsPath(ID)]
	if !ok || mountOutput.Type != MountType {
		fmt.Printf("4\n")
		return nil, microerror.Maskf(notFoundError, "PKI backend for ID '%s'", ID)
	}

	fmt.Printf("5\n")

	return mountOutput, nil
}

// GetCACertificate returns the public key of the root CA of the PKI backend
// associated to the given ID, if any.
func (p *VaultPKI) GetCACertificate(ID string) (string, error) {
	fmt.Printf("start VaultPKI.GetCACertificate\n")
	defer fmt.Printf("end VaultPKI.GetCACertificate\n")

	fmt.Printf("1\n")

	secret, err := p.vaultClient.Logical().Read(key.ReadCAPath(ID))
	if IsNoVaultHandlerDefined(err) {
		fmt.Printf("2\n")
		return "", microerror.Maskf(notFoundError, "root CA for ID '%s'", ID)
	} else if err != nil {
		fmt.Printf("3\n")
		return "", microerror.Mask(err)
	}

	// If the secret is nil, the CA has not been generated.
	if secret == nil {
		fmt.Printf("4\n")
		return "", microerror.Maskf(notFoundError, "root CA for ID '%s'", ID)
	}

	var crt string
	{
		v, ok := secret.Data["certificate"]
		if !ok {
			return "", microerror.Maskf(executionFailedError, "certificate missing")
		}
		crt, ok = v.(string)
		if !ok {
			return "", microerror.Maskf(executionFailedError, "certificate must be string")
		}
	}
	fmt.Printf("5\n")

	return crt, nil
}
