package vaultcrt

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	currentSecret, err := toSecret(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the secret has to be deleted")

	var secretToDelete *apiv1.Secret
	if currentSecret != nil {
		secretToDelete = currentSecret
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found out if the secret has to be deleted")

	return secretToDelete, nil
}

func (r *Resource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	secretToDelete, err := toSecret(deleteState)
	if err != nil {
		return microerror.Mask(err)
	}

	if secretToDelete != nil {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleting the sercet in the Kubernetes API")

		ns := key.SecretNamespace(customObject)
		if ns == "" {
			ns = r.namespace
		}

		err = r.k8sClient.CoreV1().Secrets(ns).Delete(secretToDelete.Name, &apismetav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleted the sercet in the Kubernetes API")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "the sercet does not need to be deleted from the Kubernetes API")
	}

	return nil
}
