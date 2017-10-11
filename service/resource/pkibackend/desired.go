package pkibackend

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/flannel-operator/service/key"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "computing the desired PKI backend state")

	caState := CAState{
		IsBackendMounted: true,
		IsCAGenerated:    true,
		IsRoleCreated:    true,
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "computed the desired PKI backend state")

	return caState, nil
}
