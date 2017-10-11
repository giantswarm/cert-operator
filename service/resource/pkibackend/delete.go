package pkibackend

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/flannel-operator/service/key"
)

func (r *Resource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	currentNamespace, err := toNamespace(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredNamespace, err := toNamespace(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the namespace has to be deleted")

	var namespaceToDelete *apiv1.Namespace
	if currentNamespace != nil {
		namespaceToDelete = desiredNamespace
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found out if the namespace has to be deleted")

	return namespaceToDelete, nil
}

func (r *Resource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	namespaceToDelete, err := toNamespace(deleteState)
	if err != nil {
		return microerror.Mask(err)
	}

	if namespaceToDelete != nil {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleting the namespace in the Kubernetes API")

		err = r.k8sClient.CoreV1().Namespaces().Delete(namespaceToDelete.Name, &apismetav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "deleted the namespace in the Kubernetes API")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "the namespace does not need to be deleted from the Kubernetes API")
	}

	return nil
}

//

func (s *service) Delete(clusterID string) error {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := s.VaultClient.Sys()

	// Unmount the PKI backend, if it exists.
	mounted, err := s.IsMounted(clusterID)
	if err != nil {
		return microerror.Mask(err)
	}
	if mounted {
		err = sysBackend.Unmount(s.MountPKIPath(clusterID))
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
