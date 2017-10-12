package vaultpki

import (
	"fmt"

	"github.com/giantswarm/microerror"
	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/vaultpki/key"
)

func (p *VaultPKI) CreateBackend(ID string) error {
	k := key.MountPKIPath(ID)
	v := &vaultclient.MountInput{
		Config: vaultclient.MountConfigInput{
			MaxLeaseTTL: p.caTTL,
		},
		Description: fmt.Sprintf("PKI backend for cluster ID '%s'", ID),
		Type:        MountType,
	}

	err := p.vaultClient.Sys().Mount(k, v)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (p *VaultPKI) CreateCA(ID string) error {
	k := key.WriteCAPath(ID)
	v := map[string]interface{}{
		"common_name": key.CommonName(ID, p.commonNameFormat),
		"ttl":         p.caTTL,
	}

	_, err := p.vaultClient.Logical().Write(k, v)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
