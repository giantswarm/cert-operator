package vaultaccess

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/reconciliationcanceledcontext"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.LogCtx(ctx, "level", "debug", "message", "renewing the Vault token")
	_, err := r.vaultClient.Auth().Token().RenewSelf(0)
	if IsVaultAccess(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "vault not reachable")
		r.logger.LogCtx(ctx, "level", "debug", "message", "vault upgrade in progress")
		r.logger.LogCtx(ctx, "level", "debug", "message", "canceling reconciliation")
		reconciliationcanceledcontext.SetCanceled(ctx)
		return nil

	} else if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "level", "debug", "message", "renewed the Vault token")

	return nil
}
