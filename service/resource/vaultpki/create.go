package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	changes, err := toChanges(createChange)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, change := range changes {
		switch change {
		case BackendChange:
			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "creating the Vault PKI in the Vault API")

			err := r.vaultPKI.CreateBackend(key.ClusterID(customObject))
			if err != nil {
				return microerror.Mask(err)
			}

			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "created the Vault PKI in the Vault API")
		case CACertificateChange:
			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "creating the root CA in the Vault PKI")

			err := r.vaultPKI.CreateCA(key.ClusterID(customObject))
			if err != nil {
				return microerror.Mask(err)
			}

			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "created the root CA in the Vault PKI")
		default:
			return microerror.Maskf(unknownChangeTypeError, "change=%v", change)
		}
	}

	return nil
}
