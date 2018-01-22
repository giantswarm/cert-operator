package vaultrolev1

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultrole"

	"github.com/giantswarm/cert-operator/service/keyv2"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := keyv2.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "debug", "computing the desired role")

	TTL, err := time.ParseDuration(keyv2.RoleTTL(customObject))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	role := &vaultrole.Role{
		AllowBareDomains: keyv2.AllowBareDomains(customObject),
		AllowSubdomains:  AllowSubdomains,
		AltNames:         keyv2.AltNames(customObject),
		ID:               keyv2.ClusterID(customObject),
		Organizations:    keyv2.Organizations(customObject),
		TTL:              TTL,
	}

	r.logger.LogCtx(ctx, "debug", "computed the desired role")

	return role, nil
}
