package vaultcrt

import (
	"context"

	"github.com/giantswarm/microerror"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	apiv1alpha2 "sigs.k8s.io/cluster-api/api/v1alpha2"

	"github.com/giantswarm/cert-operator/service/controller/key"
)

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	secretToCreate, err := toSecret(createChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if secretToCreate != nil {
		customObject, err := key.ToCustomObject(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "finding cluster resource")
		cluster := &apiv1alpha2.Cluster{}
		if err := r.ctrlClient.Get(ctx, types.NamespacedName{
			Namespace: customObject.Namespace,
			Name:      key.ClusterID(customObject)},
			cluster); err != nil {
			return microerror.Maskf(notFoundError, "Could not find cluster %s in namespace %s.",
				customObject.Namespace,
				key.ClusterID(customObject))
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "creating the secret in the Kubernetes API")

		_, err = r.k8sClient.CoreV1().Secrets(customObject.GetNamespace()).Create(ctx, secretToCreate, metav1.CreateOptions{})
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
