package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
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

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the PKI backend has to be created")

	var vaultPKIStateToCreate VaultPKIState
	if !currentVaultPKIState.BackendExists || !currentVaultPKIState.CAExists || !currentVaultPKIState.IsPolicyCreated || !currentVaultPKIState.IsRoleCreated {
		vaultPKIStateToCreate = desiredVaultPKIState
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found out if the PKI backend has to be created")

	return vaultPKIStateToCreate, nil
}

func (r *Resource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	vaultPKIStateToCreate, err := toVaultPKIState(createState)
	if err != nil {
		return microerror.Mask(err)
	}

	if !vaultPKIStateToCreate.BackendExists || !vaultPKIStateToCreate.CAExists || !vaultPKIStateToCreate.IsRoleCreated {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "creating the PKI backend in the Kubernetes API")

		if !vaultPKIStateToCreate.BackendExists {
			err := r.vaultPKI.CreateBackend(key.ClusterID(customObject))
			if err != nil {
				return microerror.Mask(err)
			}
		}

		if !vaultPKIStateToCreate.CAExists {
			err := r.vaultPKI.CreateCA(key.ClusterID(customObject))
			if err != nil {
				return microerror.Mask(err)
			}
		}

		if !vaultPKIStateToCreate.IsPolicyCreated {
			k := key.VaultPolicyName(customObject)
			v := `
				path "pki-` + key.ClusterID(customObject) + `/issue/role-` + key.ClusterID(customObject) + `" {
					policy = "write"
				}
			`

			err := r.vaultClient.Sys().PutPolicy(k, v)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		if !vaultPKIStateToCreate.IsRoleCreated {
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
