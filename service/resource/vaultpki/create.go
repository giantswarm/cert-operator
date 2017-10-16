package vaultpki

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	currentVaultPKIState, err := toVaultPKIState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the Vault PKI has to be created")

	fmt.Printf("%#v\n", currentVaultPKIState)

	fmt.Printf("1\n")

	var vaultPKIStateToCreate VaultPKIState
	if currentVaultPKIState.BackendExists {
		fmt.Printf("2\n")
		vaultPKIStateToCreate.BackendExists = currentVaultPKIState.BackendExists
	}
	fmt.Printf("3\n")
	if currentVaultPKIState.CAExists {
		fmt.Printf("4\n")
		vaultPKIStateToCreate.CAExists = currentVaultPKIState.CAExists
	}

	fmt.Printf("5\n")

	fmt.Printf("%#v\n", vaultPKIStateToCreate)

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found out if the Vault PKI has to be created")

	return vaultPKIStateToCreate, nil
}

func (r *Resource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	vaultPKIStateToCreate, err := toVaultPKIState(createState)
	if err != nil {
		return microerror.Mask(err)
	}

	fmt.Printf("%#v\n", vaultPKIStateToCreate)

	if !vaultPKIStateToCreate.BackendExists {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "creating the Vault PKI in the Vault API")
		err := r.vaultPKI.CreateBackend(key.ClusterID(customObject))
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "created the Vault PKI in the Vault API")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "the Vault PKI does not need to be created in the Vault API")
	}

	if !vaultPKIStateToCreate.CAExists {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "creating the root CA in the Vault PKI")
		err := r.vaultPKI.CreateCA(key.ClusterID(customObject))
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "created the root CA in the Vault PKI")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "the root CA does not need to be created in the Vault PKI")
	}

	return nil
}
