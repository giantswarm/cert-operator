package vaultpki

import (
	"context"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "looking for the PKI backend state in the Vault API")

	var vaultPKIState VaultPKIState
	{
		vaultPKIState.BackendExists, err = r.vaultPKI.BackendExists(key.ClusterID(customObject))
		if err != nil {
			return false, microerror.Mask(err)
		}

		vaultPKIState.CAExists, err = r.vaultPKI.CAExists(key.ClusterID(customObject))
		if err != nil {
			return false, microerror.Mask(err)
		}

		vaultPKIState.IsPolicyCreated, err = r.isPolicyCreated(customObject)
		if err != nil {
			return false, microerror.Mask(err)
		}

		vaultPKIState.IsRoleCreated, err = r.isRoleCreated(customObject)
		if err != nil {
			return false, microerror.Mask(err)
		}
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found the PKI backend state in the Vault API")

	return vaultPKIState, nil
}

func (r *Resource) isPolicyCreated(customObject certificatetpr.CustomObject) (bool, error) {
	// Get the system backend for policy operations.
	sysBackend := r.vaultClient.Sys()

	// Check if the policy is already there.
	policies, err := sysBackend.ListPolicies()
	if err != nil {
		return false, microerror.Mask(err)
	}
	for _, p := range policies {
		if p == key.VaultPolicyName(customObject) {
			return true, nil
		}
	}

	return false, nil
}

func (r *Resource) isRoleCreated(customObject certificatetpr.CustomObject) (bool, error) {
	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's PKI backend.
	logicalBackend := r.vaultClient.Logical()

	// Check if a PKI for the given cluster ID exists.
	secret, err := logicalBackend.List(key.VaultListRolesPath(customObject))
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}

	// In case there is not a single role for this PKI backend, secret is nil.
	if secret == nil {
		return false, nil
	}

	// When listing roles a list of role names is returned. Here we iterate over
	// this list and if we find the desired role name, it means the role has
	// already been created.
	if keys, ok := secret.Data["keys"]; ok {
		if list, ok := keys.([]interface{}); ok {
			for _, k := range list {
				if str, ok := k.(string); ok && str == key.VaultRoleName(customObject) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
