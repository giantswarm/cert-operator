package vaultcrt

import (
	"context"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/certs"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
	apiv1 "k8s.io/api/core/v1"

	"github.com/giantswarm/cert-operator/service/certconfig/v2/key"
)

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	secretToUpdate, err := toSecret(updateChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if secretToUpdate != nil {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "updating the secret in the Kubernetes API")

		_, err := r.k8sClient.CoreV1().Secrets(r.namespace).Update(secretToUpdate)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "updated the secret in the Kubernetes API")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "the secret does not need to be updated in the Kubernetes API")
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

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the secret has to be updated")

	var secretToUpdate *apiv1.Secret
	{
		TTL, err := time.ParseDuration(key.CrtTTL(customObject))
		if err != nil {
			return false, microerror.Mask(err)
		}

		renew, err := r.shouldCertBeRenewed(customObject, currentSecret, desiredSecret, TTL, r.expirationThreshold)
		if IsMissingAnnotation(err) {
			// fall through
		} else if err != nil {
			return nil, microerror.Mask(err)
		}

		if renew {
			ca, crt, k, err := r.issueCertificate(customObject)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			secretToUpdate = desiredSecret
			secretToUpdate.StringData[key.CAID] = ca
			secretToUpdate.StringData[key.CrtID] = crt
			secretToUpdate.StringData[key.KeyID] = k
		}
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found out if the secret has to be updated")

	return secretToUpdate, nil
}

func (r *Resource) shouldCertBeRenewed(customObject v1alpha1.CertConfig, currentSecret, desiredSecret *apiv1.Secret, TTL, threshold time.Duration) (bool, error) {
	// Check if there are annotations at all.
	{
		if currentSecret == nil {
			return false, microerror.Maskf(missingAnnotationError, "current secret")
		}
		if currentSecret.Annotations == nil {
			return false, microerror.Maskf(missingAnnotationError, "current secret")
		}
		if desiredSecret == nil {
			return false, microerror.Maskf(missingAnnotationError, "desired secret")
		}
		if desiredSecret.Annotations == nil {
			return false, microerror.Maskf(missingAnnotationError, "desired secret")
		}
	}

	// Check if the cert configs ask to disable regeneration.
	{
		// TODO remove this hack once all cert configs are updated with the correct
		// value for DisableRegeneration.
		if customObject.Spec.Cert.ClusterComponent == string(certs.ServiceAccountCert) {
			return false, nil
		}
		if customObject.Spec.Cert.DisableRegeneration {
			return false, nil
		}
	}

	// Check the update timestamp annotation.
	{
		a, ok := currentSecret.Annotations[UpdateTimestampAnnotation]
		if !ok {
			return false, microerror.Maskf(missingAnnotationError, "current secret")
		}

		t, err := time.ParseInLocation(UpdateTimestampLayout, a, time.UTC)
		if err != nil {
			return false, microerror.Mask(err)
		}

		if t.Add(TTL).Add(-threshold).Before(r.currentTimeFactory()) {
			return true, nil
		}
	}

	// Check the config hash annotation.
	{
		c, ok := currentSecret.Annotations[ConfigHashAnnotation]
		if !ok {
			return true, nil
		}
		d, ok := desiredSecret.Annotations[ConfigHashAnnotation]
		if !ok {
			return false, microerror.Maskf(missingAnnotationError, "desired secret")
		}

		if c != d {
			return true, nil
		}
	}

	return false, nil
}
