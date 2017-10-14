package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework/context/deletionallowedcontext"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
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
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the Vault PKI has to be deleted")

		if currentVaultPKIState.BackendExists || currentVaultPKIState.CAExists {
			vaultPKIStateToDelete = desiredVaultPKIState
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found out if the Vault PKI has to be deleted")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "not computing delete state because Vault PKIs are not allowed to be deleted")
	}

	return vaultPKIStateToDelete, nil
}

func (r *Resource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	vaultPKIStateToDelete, err := toVaultPKIState(deleteState)
	if err != nil {
		return microerror.Mask(err)
	}

	if vaultPKIStateToDelete.BackendExists || vaultPKIStateToDelete.CAExists {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleting the Vault PKI in the Vault API")

		err := r.vaultPKI.DeleteBackend(key.ClusterID(customObject))
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleted the Vault PKI in the Vault API")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "the Vault PKI does not need to be deleted from the Vault API")
	}

	return nil
}
