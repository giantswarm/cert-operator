package pkibackend

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/flannel-operator/service/key"
)

func (r *Resource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	currentCAState, err := toCAState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredCAState, err := toCAState(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the PKI backend has to be created")

	var caStateToCreate CAState
	if !currentCAState.IsBackendMounted || !currentCAState.IsCAGenerated || !currentCAState.IsRoleCreated {
		caStateToCreate = desiredCAState
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found out if the PKI backend has to be created")

	return caStateToCreate, nil
}

func (r *Resource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	caStateToCreate, err := toCAState(createState)
	if err != nil {
		return microerror.Mask(err)
	}

	if !caStateToCreate.IsBackendMounted || !caStateToCreate.IsCAGenerated || !caStateToCreate.IsRoleCreated {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "creating the PKI backend in the Kubernetes API")

		if !caStateToCreate.IsBackendMounted {
			k := key.VaultMountPKIPath(customObject)
			v := &vaultclient.MountInput{
				Config: vaultclient.MountConfigInput{
					MaxLeaseTTL: r.caTTL,
				},
				Description: fmt.Sprintf("PKI backend for cluster ID '%s'", key.ClusterID(customObject)),
				Type:        VaultMountType,
			}

			err := r.vaultClient.Sys().Mount(k, v)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		if !caStateToCreate.IsCAGenerated {
			k := key.VaultWriteCAPath(customObject)
			v := map[string]interface{}{
				"common_name": key.VaultCommonName(customObject, r.commonNameFormat),
				"ttl":         r.caTTL,
			}

			_, err := r.vaultClient.Logical().Write(k, v)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		if !caStateToCreate.IsRoleCreated {
			k := key.VaultWriteRolePath(customObject)
			v := map[string]interface{}{
				"allow_bare_domains": key.VaultAllowBareDomains(customObject),
				"allow_subdomains":   VaultAllowSubDomains,
				"allowed_domains":    key.VaultAllowedDomains(customObject, r.commonNameFormat),
				"ttl":                r.caTTL,
			}

			_, err := r.vaultClient.Logical().Write(k, v)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "created the PKI backend in the Kubernetes API")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "the PKI backend does not need to be created in the Kubernetes API")
	}

	return nil
}
