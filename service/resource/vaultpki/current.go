package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultpki"

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
		vaultPKIState.Backend, err = r.vaultPKI.GetBackend(key.ClusterID(customObject))
		if vaultpki.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return false, microerror.Mask(err)
		}

		vaultPKIState.CACertificate, err = r.vaultPKI.GetCACertificate(key.ClusterID(customObject))
		if vaultpki.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return false, microerror.Mask(err)
		}
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found the Vault PKI in the Vault API")

	return vaultPKIState, nil
}
