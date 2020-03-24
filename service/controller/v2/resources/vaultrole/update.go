package vaultrole

import (
	"context"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/resource/crud"
	"github.com/giantswarm/vaultrole"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	roleToUpdate, err := toRole(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if roleToUpdate != nil {
		r.logger.LogCtx(ctx, "debug", "updating the role in the Vault API") // nolint: errcheck

		c := vaultrole.UpdateConfig{
			AllowBareDomains: roleToUpdate.AllowBareDomains,
			AllowSubdomains:  roleToUpdate.AllowSubdomains,
			AltNames:         roleToUpdate.AltNames,
			ID:               roleToUpdate.ID,
			Organizations:    roleToUpdate.Organizations,
			TTL:              roleToUpdate.TTL.String(),
		}
		err = r.vaultRole.Update(c)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "debug", "updated the role in the Vault API") // nolint: errcheck
	} else {
		r.logger.LogCtx(ctx, "debug", "the role does not need to be updated in the Vault API") // nolint: errcheck
	}

	return nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	create, err := r.newCreateChange(ctx, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	update, err := r.newUpdateChange(ctx, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := crud.NewPatch()
	patch.SetCreateChange(create)
	patch.SetUpdateChange(update)

	return patch, nil
}

func (r *Resource) newUpdateChange(ctx context.Context, currentState, desiredState interface{}) (interface{}, error) {
	currentRole, err := toRole(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredRole, err := toRole(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "debug", "finding out if the role has to be updated") // nolint: errcheck

	var roleToUpdate *vaultrole.Role
	if !reflect.DeepEqual(currentRole, desiredRole) {
		roleToUpdate = desiredRole
	}

	r.logger.LogCtx(ctx, "debug", "found out if the role has to be updated") // nolint: errcheck

	return roleToUpdate, nil
}
