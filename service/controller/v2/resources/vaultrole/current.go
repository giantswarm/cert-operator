package vaultrole

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultrole"

	"github.com/giantswarm/cert-operator/service/certconfig/v2/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "debug", "looking for the role in the Vault API")

	var role *vaultrole.Role
	{
		c := vaultrole.SearchConfig{
			ID:            key.ClusterID(customObject),
			Organizations: key.Organizations(customObject),
		}
		result, err := r.vaultRole.Search(c)
		if vaultrole.IsNotFound(err) {
			r.logger.LogCtx(ctx, "debug", "did not find the role in the Vault API")
			// fall through
		} else if err != nil {
			return nil, microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "debug", "found the role in the Vault API")
			role = &result
		}
	}

	return role, nil
}
