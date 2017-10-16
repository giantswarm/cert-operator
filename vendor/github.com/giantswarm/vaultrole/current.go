package vaultrole

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/vaultrole/key"
)

func (r *VaultRole) Exists(config ExistsConfig) (bool, error) {
	// Check if a PKI for the given cluster ID exists.
	secret, err := r.vaultClient.Logical().List(key.ListRolesPath(config.ID))
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
				if str, ok := k.(string); ok && str == key.RoleName(config.ID, config.Organizations) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
