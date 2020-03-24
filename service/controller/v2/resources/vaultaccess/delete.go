package vaultaccess

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/controller/context/reconciliationcanceledcontext"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	_, err := r.vaultClient.Auth().Token().LookupSelf()
	if IsVaultAccess(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "vault not reachable") // nolint: errcheck
		r.logger.LogCtx(ctx, "level", "debug", "message", "vault upgrade in progress") // nolint: errcheck
		r.logger.LogCtx(ctx, "level", "debug", "message", "canceling reconciliation") // nolint: errcheck
		reconciliationcanceledcontext.SetCanceled(ctx)
		return nil

	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
