package vaultpki

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/vaultpki/key"
)

func (p *VaultPKI) DeleteBackend(ID string) error {
	k := key.MountPKIPath(ID)
	err := p.vaultClient.Sys().Unmount(k)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
