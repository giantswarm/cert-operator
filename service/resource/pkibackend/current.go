package pkibackend

import (
	"context"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/flannel-operator/service/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "looking for the PKI backend state in the Vault API")

	var caState CAState
	{
		caState.isBackendMounted, err = r.isBackendMounted(customObject)
		if err != nil {
			return false, microerror.Mask(err)
		}

		catState.IsCAGenerated, err = r.IsCAGenerated(customObject)
		if err != nil {
			return false, microerror.Mask(err)
		}

		caState.IsRoleCreated, err = r.IsRoleCreated(customObject)
		if err != nil {
			return false, microerror.Mask(err)
		}
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found the PKI backend state in the Vault API")

	return validPKIBackend, nil
}

func (r *Resource) isBackendMounted(customObject certificatetpr.CustomObject) (bool, error) {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := r.vaultClient.Sys()

	// Check if a PKI for the given cluster ID exists.
	mounts, err := sysBackend.ListMounts()
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}
	mountOutput, ok := mounts[key.VaultListMountsPath(customObject)+"/"]
	if !ok || mountOutput.Type != "pki" {
		return false, nil
	}

	return true, nil
}

func (r *Resource) isCAGenerated(customObject certificatetpr.CustomObject) (bool, error) {
	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's PKI backend.
	logicalBackend := r.vaultClient.Logical()

	// Check if a root CA for the given cluster ID exists.
	secret, err := logicalBackend.Read(key.VaultReadCAPath(customObject))
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}

	// If the secret is nil, the CA has not been generated.
	if secret == nil {
		return false, nil
	}

	certificate, ok := secret.Data["certificate"]
	if ok && certificate == "" {
		return false, nil
	}
	err, ok = secret.Data["error"]
	if ok && err != "" {
		return false, nil
	}

	return true, nil
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
