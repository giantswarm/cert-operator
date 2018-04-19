package vaultcrt

import (
	"context"
	"fmt"
	"strings"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/controller/context/reconciliationcanceledcontext"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cert-operator/service/controller/v2/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "looking for the secret in the Kubernetes API")

	var secret *corev1.Secret
	{
		manifest, err := r.k8sClient.Core().Secrets(r.namespace).Get(key.SecretName(customObject), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", "did not find the secret in the Kubernetes API")
			// fall through
		} else if err != nil {
			return nil, microerror.Mask(err)
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", "found the secret in the Kubernetes API")
			secret = manifest
			r.updateVersionBundleVersionGauge(ctx, customObject, versionBundleVersionGauge, secret)
		}
	}

	// In case a cluster deletion happens, we want to delete all secrets holding
	// certificates. We still need the certificates for draining nodes on KVM
	// though. So as long as pods are there we delay the deletion of the secrets
	// here in order to still use them in the kvm-operator. The impact of this for
	// AWS and Azure is zero, because when listing on namespaces that do not exist
	// we get an empty list and thus do nothing here. For KVM, as soon as the
	// draining was done and the pods got removed we get an empty list here after
	// the delete event got replayed. Then we just remove the secrets as usual.
	if key.IsInDeletionState(customObject) {
		n := key.ClusterNamespace(customObject)
		list, err := r.k8sClient.CoreV1().Pods(n).List(metav1.ListOptions{})
		if err != nil {
			return nil, microerror.Mask(err)
		}
		if len(list.Items) != 0 {
			r.logger.LogCtx(ctx, "level", "debug", "message", "cannot finish deletion of the secret due to existing pods")
			reconciliationcanceledcontext.SetCanceled(ctx)
			r.logger.LogCtx(ctx, "level", "debug", "message", "canceling resource for custom object")

			return nil, nil
		}
	}

	return secret, nil
}

func (r *Resource) updateVersionBundleVersionGauge(ctx context.Context, customObject v1alpha1.CertConfig, gauge *prometheus.GaugeVec, secret *corev1.Secret) {
	version, ok := secret.Annotations[VersionBundleVersionAnnotation]
	if !ok {
		r.logger.LogCtx(ctx, "level", "warning", "message", fmt.Sprintf("cannot update current version bundle version metric: annotation '%s' must not be empty", VersionBundleVersionAnnotation))
		return
	}

	split := strings.Split(version, ".")
	if len(split) != 3 {
		r.logger.LogCtx(ctx, "level", "warning", "message", fmt.Sprintf("cannot update current version bundle version metric: invalid version format, expected '<major>.<minor>.<patch>', got '%s'", version))
		return
	}

	major := split[0]
	minor := split[1]
	patch := split[2]

	gauge.WithLabelValues(major, minor, patch).Set(1)
}
