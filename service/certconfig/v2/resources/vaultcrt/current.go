package vaultcrt

import (
	"context"
	"fmt"
	"strings"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cert-operator/service/certconfig/v2/key"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "looking for the secret in the Kubernetes API")

	var secret *apiv1.Secret
	{
		manifest, err := r.k8sClient.Core().Secrets(r.namespace).Get(key.SecretName(customObject), apismetav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "did not find the secret in the Kubernetes API")
			// fall through
		} else if err != nil {
			return nil, microerror.Mask(err)
		} else {
			r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found the secret in the Kubernetes API")
			secret = manifest
			r.updateVersionBundleVersionGauge(customObject, versionBundleVersionGauge, secret)
		}
	}

	return secret, nil
}

func (r *Resource) updateVersionBundleVersionGauge(customObject v1alpha1.CertConfig, gauge *prometheus.GaugeVec, secret *apiv1.Secret) {
	version, ok := secret.Annotations[VersionBundleVersionAnnotation]
	if !ok {
		r.logger.Log("cluster", key.ClusterID(customObject), "warning", fmt.Sprintf("cannot update current version bundle version metric: annotation '%s' must not be empty", VersionBundleVersionAnnotation))
		return
	}

	split := strings.Split(version, ".")
	if len(split) != 3 {
		r.logger.Log("cluster", key.ClusterID(customObject), "warning", fmt.Sprintf("cannot update current version bundle version metric: invalid version format, expected '<major>.<minor>.<patch>', got '%s'", version))
		return
	}

	major := split[0]
	minor := split[1]
	patch := split[2]

	gauge.WithLabelValues(major, minor, patch).Set(1)
}
