package cleaning

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cert-operator/service/controller/v2/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	clusterID := key.ClusterID(customObject)

	_, err = r.k8sClient.CoreV1().Namespaces().Get(clusterID, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("cluster namespace %#q does not exist in CP", clusterID))
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting %#q certconfig", customObject.GetName()))

		err = r.g8sClient.CoreV1alpha1().CertConfigs(customObject.GetNamespace()).Delete(customObject.Name, &metav1.DeleteOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleted %#q certconfig", customObject.GetName()))
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
