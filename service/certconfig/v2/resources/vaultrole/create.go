package vaultrole

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultrole"
)

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	roleToCreate, err := toRole(createChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if roleToCreate != nil {
		r.logger.LogCtx(ctx, "debug", "creating the role in the Vault API")

		c := vaultrole.CreateConfig{
			AllowBareDomains: roleToCreate.AllowBareDomains,
			AllowSubdomains:  roleToCreate.AllowSubdomains,
			AltNames:         roleToCreate.AltNames,
			ID:               roleToCreate.ID,
			Organizations:    roleToCreate.Organizations,
			TTL:              roleToCreate.TTL.String(),
		}
		err = r.vaultRole.Create(c)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "debug", "created the role in the Vault API")
	} else {
		r.logger.LogCtx(ctx, "debug", "the role does not need to be created in the Vault API")
	}

	return nil
}

func (r *Resource) newCreateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentRole, err := toRole(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredRole, err := toRole(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "debug", "finding out if the role has to be created")

	var roleToCreate *vaultrole.Role
	if currentRole == nil {
		roleToCreate = desiredRole
	}

	r.logger.LogCtx(ctx, "debug", "found out if the role has to be created")

	return roleToCreate, nil
}
