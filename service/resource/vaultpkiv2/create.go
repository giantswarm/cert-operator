package vaultpkiv2

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/service/keyv2"
)

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	customObject, err := keyv2.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	vaultPKIStateToCreate, err := toVaultPKIState(createChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if vaultPKIStateToCreate.Backend != nil {
		r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "creating the Vault PKI in the Vault API")

		err := r.vaultPKI.CreateBackend(keyv2.ClusterID(customObject))
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "created the Vault PKI in the Vault API")
	} else {
		r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "the Vault PKI does not need to be created in the Vault API")
	}

	if vaultPKIStateToCreate.CACertificate != "" {
		r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "creating the root CA in the Vault PKI")

		err := r.vaultPKI.CreateCA(keyv2.ClusterID(customObject))
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "created the root CA in the Vault PKI")
	} else {
		r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "the root CA does not need to be created in the Vault PKI")
	}

	return nil
}

func (r *Resource) newCreateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := keyv2.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	currentVaultPKIState, err := toVaultPKIState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredVaultPKIState, err := toVaultPKIState(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "finding out if the Vault PKI has to be created")

	var vaultPKIStateToCreate VaultPKIState
	if currentVaultPKIState.Backend == nil {
		vaultPKIStateToCreate.Backend = desiredVaultPKIState.Backend
	}
	if currentVaultPKIState.CACertificate == "" {
		vaultPKIStateToCreate.CACertificate = desiredVaultPKIState.CACertificate
	}

	r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "found out if the Vault PKI has to be created")

	return vaultPKIStateToCreate, nil
}
