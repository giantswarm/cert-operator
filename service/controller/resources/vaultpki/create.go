package vaultpki

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/service/controller/key"
)

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	vaultPKIStateToCreate, err := toVaultPKIState(createChange)
	if err != nil {
		return microerror.Mask(err)
	}

	ids := []string{
		key.ClusterID(customObject),
		fmt.Sprintf("%s-etcd", key.ClusterID(customObject)),
	}

	if vaultPKIStateToCreate.Backend != nil {
		for _, id := range ids {
			r.logger.Debugf(ctx, "creating the Vault PKI %s in the Vault API", id)

			err := r.vaultPKI.CreateBackend(id)
			if err != nil {
				return microerror.Mask(err)
			}

			r.logger.Debugf(ctx, "created the Vault PKI %s in the Vault API", id)
		}
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "the Vault PKI does not need to be created in the Vault API")
	}

	if vaultPKIStateToCreate.CACertificate != "" {
		for _, id := range ids {
			r.logger.Debugf(ctx, "creating the root CA for %s in the Vault PKI", id)

			_, err := r.vaultPKI.CreateCA(id)
			if err != nil {
				return microerror.Mask(err)
			}

			r.logger.Debugf(ctx, "created the root CA for %s in the Vault PKI", id)
		}
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "the root CA does not need to be created in the Vault PKI")
	}

	return nil
}

func (r *Resource) newCreateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentVaultPKIState, err := toVaultPKIState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredVaultPKIState, err := toVaultPKIState(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "finding out if the Vault PKI has to be created")

	var vaultPKIStateToCreate VaultPKIState
	if currentVaultPKIState.Backend == nil {
		vaultPKIStateToCreate.Backend = desiredVaultPKIState.Backend
	}
	if currentVaultPKIState.CACertificate == "" {
		vaultPKIStateToCreate.CACertificate = desiredVaultPKIState.CACertificate
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "found out if the Vault PKI has to be created")

	return vaultPKIStateToCreate, nil
}
