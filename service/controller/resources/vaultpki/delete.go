package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v4/pkg/resource/crud"

	"github.com/giantswarm/cert-operator/service/controller/key"
)

func (r *Resource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	delete, err := r.newDeleteChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := crud.NewPatch()
	patch.SetDeleteChange(delete)

	return patch, nil
}

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	vaultPKIStateToDelete, err := toVaultPKIState(deleteChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if vaultPKIStateToDelete.Backend != nil || vaultPKIStateToDelete.CACertificate != "" {
		ids := key.PKIIdsForCluster(key.ClusterID(customObject))

		for _, id := range ids {
			r.logger.Debugf(ctx, "deleting the Vault PKI %s in the Vault API", id)

			err := r.vaultPKI.DeleteBackend(id)
			if err != nil {
				return microerror.Mask(err)
			}

			r.logger.Debugf(ctx, "deleted the Vault PKI %s in the Vault API", id)
		}
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "the Vault PKIs does not need to be deleted from the Vault API")
	}

	return nil
}

func (r *Resource) newDeleteChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	_, err := toVaultPKIState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	_, err = toVaultPKIState(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// We do not delete tenant cluster PKI when a CertConfig is deleted.
	var vaultPKIStateToDelete VaultPKIState

	return vaultPKIStateToDelete, nil
}
