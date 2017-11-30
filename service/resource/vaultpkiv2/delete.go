package vaultpkiv1

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/context/deletionallowedcontext"

	"github.com/giantswarm/cert-operator/service/keyv1"
)

func (r *Resource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	delete, err := r.newDeleteChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := framework.NewPatch()
	patch.SetDeleteChange(delete)

	return patch, nil
}

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	customObject, err := keyv1.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	vaultPKIStateToDelete, err := toVaultPKIState(deleteChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if vaultPKIStateToDelete.Backend != nil || vaultPKIStateToDelete.CACertificate != "" {
		r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "deleting the Vault PKI in the Vault API")

		err := r.vaultPKI.DeleteBackend(keyv1.ClusterID(customObject))
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "deleted the Vault PKI in the Vault API")
	} else {
		r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "the Vault PKI does not need to be deleted from the Vault API")
	}

	return nil
}

func (r *Resource) newDeleteChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := keyv1.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	currentVaultPKIState, err := toVaultPKIState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredVaultPKIState, err := toVaultPKIState(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var vaultPKIStateToDelete VaultPKIState
	if deletionallowedcontext.IsDeletionAllowed(ctx) {
		r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "finding out if the Vault PKI has to be deleted")

		if currentVaultPKIState.Backend == nil {
			vaultPKIStateToDelete.Backend = desiredVaultPKIState.Backend
		}
		if currentVaultPKIState.CACertificate == "" {
			vaultPKIStateToDelete.CACertificate = desiredVaultPKIState.CACertificate
		}

		r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "found out if the Vault PKI has to be deleted")
	} else {
		r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "not computing delete state because Vault PKIs are not allowed to be deleted")
	}

	return vaultPKIStateToDelete, nil
}
