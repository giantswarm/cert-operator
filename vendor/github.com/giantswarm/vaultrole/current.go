package vaultrole

import (
	"fmt"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/vaultrole/key"
)

func (r *VaultRole) Exists(config ExistsConfig) (bool, error) {
	fmt.Printf("start VaultRole.Exists\n")
	defer fmt.Printf("end VaultRole.Exists\n")

	fmt.Printf("key.RoleName: %#v\n", key.RoleName(config.ID, config.Organizations))

	fmt.Printf("1\n")

	// Check if a PKI for the given cluster ID exists.
	secret, err := r.vaultClient.Logical().List(key.ListRolesPath(config.ID))
	if IsNoVaultHandlerDefined(err) {
		fmt.Printf("2\n")
		return false, nil
	} else if err != nil {
		fmt.Printf("3\n")
		return false, microerror.Mask(err)
	}

	// In case there is not a single role for this PKI backend, secret is nil.
	if secret == nil {
		fmt.Printf("4\n")
		return false, nil
	}

	// When listing roles a list of role names is returned. Here we iterate over
	// this list and if we find the desired role name, it means the role has
	// already been created.
	if keys, ok := secret.Data["keys"]; ok {
		fmt.Printf("5\n")
		if list, ok := keys.([]interface{}); ok {
			for _, k := range list {
				fmt.Printf("k: %#v\n", k)
				if str, ok := k.(string); ok && str == key.RoleName(config.ID, config.Organizations) {
					fmt.Printf("6\n")
					return true, nil
				}
			}
		}
	}
	fmt.Printf("7\n")

	return false, nil
}
