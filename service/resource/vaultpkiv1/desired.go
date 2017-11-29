package vaultpkiv1

import (
	"context"

	"github.com/giantswarm/cert-operator/service/keyv1"
	"github.com/giantswarm/microerror"
	vaultapi "github.com/hashicorp/vault/api"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := keyv1.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "computing the desired Vault PKI")

	// NOTE that we only define a sparse desired state. This is good enough
	// because we only need a non-zero-value desired state to do the proper
	// reconciliation. It is also that in case of the CA certificate we cannot
	// just predict and define it here, because this is the resonsibility of the
	// actual issuer backend, e.g. Vault.
	var vaultPKIState VaultPKIState
	{
		vaultPKIState.Backend = &vaultapi.MountOutput{
			Type: "pki",
		}

		vaultPKIState.CACertificate = "placeholder"
	}

	r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "computed the desired Vault PKI")

	return vaultPKIState, nil
}
