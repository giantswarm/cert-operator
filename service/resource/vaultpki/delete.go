package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	changes, err := toChanges(deleteChange)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, change := range changes {
		switch change {
		case BackendChange, CACertificateChange:
			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleting the Vault PKI in the Vault API")

			err := r.vaultPKI.DeleteBackend(key.ClusterID(customObject))
			if err != nil {
				return microerror.Mask(err)
			}

			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleted the Vault PKI in the Vault API")

			// Return fast. After removing the backend there is nothing more to change.
			return nil
		default:
			return microerror.Maskf(unknownChangeTypeError, "change=%v", change)
		}
	}

	return nil
}
