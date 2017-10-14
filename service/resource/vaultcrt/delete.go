package vaultcrt

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/context/deletionallowedcontext"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	currentSecret, err := toSecret(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredSecret, err := toSecret(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var secretToDelete Secret
	if deletionallowedcontext.IsDeletionAllowed(ctx) {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the PKI backend has to be deleted")

		if currentSecret.BackendExists || currentSecret.CAExists || currentSecret.IsPolicyCreated || currentSecret.IsRoleCreated {
			secretToDelete = desiredSecret
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found out if the PKI backend has to be deleted")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "not computing delete state because PKI backends are not allowed to be deleted")
	}

	return secretToDelete, nil
}

func (r *Resource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	secretToDelete, err := toSecret(deleteState)
	if err != nil {
		return microerror.Mask(err)
	}

	if secretToDelete.BackendExists || secretToDelete.CAExists || secretToDelete.IsPolicyCreated || secretToDelete.IsRoleCreated {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleting the PKI backend in the Vault API")

		if secretToDelete.BackendExists || secretToDelete.CAExists {
			err := r.vaultPKI.DeleteBackend(key.ClusterID(customObject))
			if err != nil {
				return microerror.Mask(err)
			}
		}

		if secretToDelete.IsPolicyCreated {
			// TODO
		}

		if secretToDelete.IsRoleCreated {
			k := key.VaultPolicyName(customObject)
			err := r.vaultClient.Sys().DeletePolicy(k)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleted the PKI backend in the Vault API")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "the PKI backend does not need to be deleted from the Vault API")
	}

	return nil
}
