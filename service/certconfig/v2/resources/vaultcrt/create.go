package vaultcrt

import (
	"context"

	"github.com/giantswarm/microerror"
	apiv1 "k8s.io/api/core/v1"

	"github.com/giantswarm/cert-operator/service/certconfig/v2/key"
)

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	secretToCreate, err := toSecret(createChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if secretToCreate != nil {
		r.logger.LogCtx(ctx, "level", "debug", "message", "creating the secret in the Kubernetes API")

		_, err := r.k8sClient.CoreV1().Secrets(r.namespace).Create(secretToCreate)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "created the secret in the Kubernetes API")
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "the secret does not need to be created in the Kubernetes API")
	}

	return nil
}

func (r *Resource) newCreateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	currentSecret, err := toSecret(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredSecret, err := toSecret(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "finding out if the secret has to be created")

	var secretToCreate *apiv1.Secret
	if currentSecret == nil {
		ca, crt, k, err := r.issueCertificate(customObject)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		secretToCreate = desiredSecret
		secretToCreate.StringData[key.CAID] = ca
		secretToCreate.StringData[key.CrtID] = crt
		secretToCreate.StringData[key.KeyID] = k
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "found out if the secret has to be created")

	return secretToCreate, nil
}
