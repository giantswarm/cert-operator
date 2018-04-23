package vaultcrt

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/controller"
	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	secretToDelete, err := toSecret(deleteChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if secretToDelete != nil {
		r.logger.LogCtx(ctx, "level", "debug", "message", "deleting the sercet in the Kubernetes API")

		err = r.k8sClient.CoreV1().Secrets(r.namespace).Delete(secretToDelete.Name, &apismetav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "deleted the sercet in the Kubernetes API")
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "the sercet does not need to be deleted from the Kubernetes API")
	}

	return nil
}

func (r *Resource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	delete, err := r.newDeleteChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := controller.NewPatch()
	patch.SetDeleteChange(delete)

	return patch, nil
}

func (r *Resource) newDeleteChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	currentSecret, err := toSecret(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "finding out if the secret has to be deleted")

	var secretToDelete *apiv1.Secret
	if currentSecret != nil {
		secretToDelete = currentSecret
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "found out if the secret has to be deleted")

	return secretToDelete, nil
}
