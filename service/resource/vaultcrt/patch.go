package vaultcrt

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
)

func (r *Resource) NewPatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	current, err := toSecret(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	desired, err := toSecret(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := framework.NewPatch()

	if current == nil && desired != nil {
		patch.SetCreateChange(desired)
	}

	if current != nil && desired == nil {
		patch.SetDeleteChange(current)
	}

	return patch, nil
}
