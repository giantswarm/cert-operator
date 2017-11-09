package vaultcrt

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	secretToDelete, err := toSecret(deleteChange)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleting the sercet in the Kubernetes API")

	err = r.k8sClient.CoreV1().Secrets(r.namespace).Delete(secretToDelete.Name, &apismetav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		// fall through
	} else if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleted the sercet in the Kubernetes API")

	return nil
}
