package vaultcrtv1

import (
	"context"
	"time"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
	apiv1 "k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/cert-operator/service/keyv1"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	customObject, err := keyv1.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	secretToUpdate, err := toSecret(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if secretToUpdate != nil {
		r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "updating the secret in the Kubernetes API")

		_, err := r.k8sClient.CoreV1().Secrets(r.namespace).Update(secretToUpdate)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "updated the secret in the Kubernetes API")
	} else {
		r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "the secret does not need to be updated in the Kubernetes API")
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
	customObject, err := keyv1.ToCustomObject(obj)
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

	r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "finding out if the secret has to be updated")

	var secretToUpdate *apiv1.Secret
	{
		TTL, err := time.ParseDuration(keyv1.CrtTTL(customObject))
		if err != nil {
			return false, microerror.Mask(err)
		}

		renew, err := r.shouldCertBeRenewed(currentSecret, TTL, r.expirationThreshold)
		if IsMissingAnnotation(err) {
			// fall through
		} else if err != nil {
			return nil, microerror.Mask(err)
		}

		if renew {
			secretToUpdate = desiredSecret

			err := r.ensureVaultRole(customObject)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			ca, crt, key, err := r.issueCertificate(customObject)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			secretToUpdate.StringData[certificatetpr.CA.String()] = ca
			secretToUpdate.StringData[certificatetpr.Crt.String()] = crt
			secretToUpdate.StringData[certificatetpr.Key.String()] = key
		}
	}

	r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "found out if the secret has to be updated")

	return secretToUpdate, nil
}

func (r *Resource) shouldCertBeRenewed(secret *apiv1.Secret, TTL, threshold time.Duration) (bool, error) {
	if secret == nil {
		return false, microerror.Mask(missingAnnotationError)
	}
	if secret.Annotations == nil {
		return false, microerror.Mask(missingAnnotationError)
	}
	a, ok := secret.Annotations[UpdateTimestampAnnotation]
	if !ok {
		return false, microerror.Mask(missingAnnotationError)
	}

	t, err := time.ParseInLocation(UpdateTimestampLayout, a, time.UTC)
	if err != nil {
		return false, microerror.Mask(err)
	}

	if t.Add(TTL).Add(-threshold).Before(r.currentTimeFactory()) {
		return true, nil
	}

	return false, nil
}
