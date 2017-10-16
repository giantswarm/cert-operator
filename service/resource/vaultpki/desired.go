package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "computing the desired Vault PKI")

	vaultPKIState := VaultPKIState{
		BackendMissing: false,
		CAMissing:      false,
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "computed the desired Vault PKI")

	return vaultPKIState, nil
}
