package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultpki"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}, deleted bool) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var vaultPKIState VaultPKIState

	{
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "looking for the Vault PKI in the Vault API")

		_, err := r.vaultPKI.GetBackend(key.ClusterID(customObject))
		if vaultpki.IsNotFound(err) {
			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "did not find the Vault PKI in the Vault API")
		} else if err != nil {
			return false, microerror.Mask(err)
		} else {
			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found the Vault PKI in the Vault API")

			vaultPKIState.BackendExists = true
		}
	}

	{
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "looking for the root CA in the Vault PKI")

		_, err := r.vaultPKI.GetCACertificate(key.ClusterID(customObject))
		if vaultpki.IsNotFound(err) {
			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "did not find the root CA in the Vault PKI")
		} else if err != nil {
			return false, microerror.Mask(err)
		} else {
			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found the root CA in the Vault PKI")

			vaultPKIState.CACertificateExists = true
		}
	}

	return vaultPKIState, nil
}
