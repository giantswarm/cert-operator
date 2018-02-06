package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultpki"

	"github.com/giantswarm/cert-operator/service/certconfig/v2/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := keyv2.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var vaultPKIState VaultPKIState

	{
		r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "looking for the Vault PKI in the Vault API")

		backend, err := r.vaultPKI.GetBackend(keyv2.ClusterID(customObject))
		if vaultpki.IsNotFound(err) {
			r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "did not find the Vault PKI in the Vault API")
		} else if err != nil {
			return false, microerror.Mask(err)
		} else {
			r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "found the Vault PKI in the Vault API")

			vaultPKIState.Backend = backend
		}
	}

	{
		r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "looking for the root CA in the Vault PKI")

		caCertificate, err := r.vaultPKI.GetCACertificate(keyv2.ClusterID(customObject))
		if vaultpki.IsNotFound(err) {
			r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "did not find the root CA in the Vault PKI")
		} else if err != nil {
			return false, microerror.Mask(err)
		} else {
			r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "found the root CA in the Vault PKI")

			vaultPKIState.CACertificate = caCertificate
		}
	}

	return vaultPKIState, nil
}
