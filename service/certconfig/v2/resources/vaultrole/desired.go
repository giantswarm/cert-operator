package vaultrole

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultrole"

	"github.com/giantswarm/cert-operator/service/certconfig/v2/key"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "debug", "computing the desired role")

	TTL, err := time.ParseDuration(key.RoleTTL(customObject))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	role := &vaultrole.Role{
		AllowBareDomains: key.AllowBareDomains(customObject),
		AllowSubdomains:  AllowSubdomains,
		AltNames:         key.AltNames(customObject),
		ID:               key.ClusterID(customObject),
		Organizations:    key.Organizations(customObject),
		TTL:              TTL,
	}

	r.logger.LogCtx(ctx, "debug", "computed the desired role")

	return role, nil
}
