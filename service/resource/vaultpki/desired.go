package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}, deleted bool) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return VaultPKIState{}, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "computing the desired Vault PKI")

	var vaultPKIState VaultPKIState
	if !deleted {
		vaultPKIState.BackendExists = true
		vaultPKIState.CACertificateExists = true
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "computed the desired Vault PKI")

	return vaultPKIState, nil
}
