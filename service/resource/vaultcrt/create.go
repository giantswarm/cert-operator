package vaultcrt

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
	currentSecret, err := toSecret(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredSecret, err := toSecret(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the secret has to be created")

	var secretToCreate *apiv1.Secret
	if currentSecret == nil {
		secretToCreate = desiredSecret

		exists, err := r.vaultRole.Exi
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found out if the secret has to be created")

	return secretToCreate, nil
}

func (r *Resource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	secretToCreate, err := toSecret(createState)
	if err != nil {
		return microerror.Mask(err)
	}

	if !secretToCreate.BackendExists || !secretToCreate.CAExists || !secretToCreate.IsRoleCreated {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "creating the PKI backend in the Kubernetes API")

		if !secretToCreate.BackendExists {
			err := r.vaultPKI.CreateBackend(key.ClusterID(customObject))
			if err != nil {
				return microerror.Mask(err)
			}
		}

		if !secretToCreate.CAExists {
			err := r.vaultPKI.CreateCA(key.ClusterID(customObject))
			if err != nil {
				return microerror.Mask(err)
			}
		}

		if !secretToCreate.IsPolicyCreated {
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

		if !secretToCreate.IsRoleCreated {
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
