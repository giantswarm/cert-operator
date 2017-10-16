package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "looking for the Vault PKI in the Vault API")

	var vaultPKIState VaultPKIState
	{
		backendExists, err := r.vaultPKI.BackendExists(key.ClusterID(customObject))
		if err != nil {
			return false, microerror.Mask(err)
		}
		vaultPKIState.BackendMissing = !backendExists

		caExists, err := r.vaultPKI.CAExists(key.ClusterID(customObject))
		if err != nil {
			return false, microerror.Mask(err)
		}
		vaultPKIState.CAMissing = !caExists
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found the Vault PKI in the Vault API")

	return vaultPKIState, nil
}
