package vaultpki

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
)

func (r *Resource) NewPatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	current, err := toVaultPKIState(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	desired, err := toVaultPKIState(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := framework.NewPatch()

	var create, delete []ChangeType

	if !current.BackendExists && desired.BackendExists {
		create = append(create, BackendChange)
	}
	if !current.CACertificateExists && desired.CACertificateExists {
		create = append(create, CACertificateChange)
	}

	if current.BackendExists && !desired.BackendExists {
		delete = append(delete, BackendChange)
	}
	if current.CACertificateExists && !desired.CACertificateExists {
		delete = append(delete, CACertificateChange)
	}

	if len(create) > 0 {
		patch.SetCreateChange(create)
	}
	if len(delete) > 0 {
		patch.SetDeleteChange(delete)
	}
	return patch, nil
}
