package vaultrole

import (
	"context"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/vaultrole"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	roleToUpdate, err := toRole(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if roleToUpdate != nil {
		r.logger.LogCtx(ctx, "debug", "updating the role in the Vault API")

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

		r.logger.LogCtx(ctx, "debug", "updated the role in the Vault API")
	} else {
		r.logger.LogCtx(ctx, "debug", "the role does not need to be updated in the Vault API")
	}

	return nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	create, err := r.newCreateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	update, err := r.newUpdateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := framework.NewPatch()
	patch.SetCreateChange(create)
	patch.SetUpdateChange(update)

	return patch, nil
}

func (r *Resource) newUpdateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentRole, err := toRole(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredRole, err := toRole(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "debug", "finding out if the role has to be updated")

	var roleToUpdate *vaultrole.Role
	if !reflect.DeepEqual(currentRole, desiredRole) {
		roleToUpdate = desiredRole
	}

	r.logger.LogCtx(ctx, "debug", "found out if the role has to be updated")

	return roleToUpdate, nil
}
