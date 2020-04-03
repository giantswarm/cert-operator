package vaultpki

import (
	"context"

	vaultapi "github.com/hashicorp/vault/api"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.LogCtx(ctx, "level", "debug", "message", "computing the desired Vault PKI")

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

	r.logger.LogCtx(ctx, "level", "debug", "message", "computed the desired Vault PKI")

	return vaultPKIState, nil
}
