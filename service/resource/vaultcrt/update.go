package vaultcrt

import (
	"context"
)

// NOTE update procedures are not implemented at the moment because we do not
// renew certificates yet.

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	return nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	return nil, nil
}
