package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultpki"

	"github.com/giantswarm/cert-operator/v3/service/controller/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var vaultPKIState VaultPKIState

	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "looking for the Vault PKI in the Vault API")

		backend, err := r.vaultPKI.GetBackend(key.ClusterID(customObject))
		if vaultpki.IsNotFound(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", "did not find the Vault PKI in the Vault API")
		} else if err != nil {
			return false, microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", "found the Vault PKI in the Vault API")

			vaultPKIState.Backend = backend
		}
	}

	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "looking for the root CA in the Vault PKI")

		caCertificate, err := r.vaultPKI.GetCACertificate(key.ClusterID(customObject))
		if vaultpki.IsNotFound(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", "did not find the root CA in the Vault PKI")
		} else if err != nil {
			return false, microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", "found the root CA in the Vault PKI")

			vaultPKIState.CACertificate = caCertificate.Certificate
		}
	}

	return vaultPKIState, nil
}
