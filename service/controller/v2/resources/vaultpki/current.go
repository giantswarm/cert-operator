package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultpki"

	"github.com/giantswarm/cert-operator/service/controller/v2/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var vaultPKIState VaultPKIState

	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "looking for the Vault PKI in the Vault API") // nolint: errcheck

		backend, err := r.vaultPKI.GetBackend(key.ClusterID(customObject))
		if vaultpki.IsNotFound(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", "did not find the Vault PKI in the Vault API") // nolint: errcheck
		} else if err != nil {
			return false, microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", "found the Vault PKI in the Vault API") // nolint: errcheck

			vaultPKIState.Backend = backend
		}
	}

	{
		r.logger.LogCtx(ctx, "level", "debug", "message", "looking for the root CA in the Vault PKI") // nolint: errcheck

		caCertificate, err := r.vaultPKI.GetCACertificate(key.ClusterID(customObject))
		if vaultpki.IsNotFound(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", "did not find the root CA in the Vault PKI") // nolint: errcheck
		} else if err != nil {
			return false, microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", "found the root CA in the Vault PKI") // nolint: errcheck

			vaultPKIState.CACertificate = caCertificate.Certificate
		}
	}

	return vaultPKIState, nil
}
