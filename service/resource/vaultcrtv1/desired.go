package vaultcrtv1

import (
	"context"
	"time"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/cert-operator/service/keyv1"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := keyv1.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "computing the desired secret")

	// NOTE that the actual secret content here is left blank because only the
	// issuer backend, e.g. Vault, can generate certificates. This has to be
	// considered when computing the create, delete and update state.
	secret := &apiv1.Secret{
		ObjectMeta: apismetav1.ObjectMeta{
			Name: keyv1.SecretName(customObject),
			Annotations: map[string]string{
				UpdateTimestampAnnotation:      r.currentTimeFactory().In(time.UTC).Format(UpdateTimestampLayout),
				VersionBundleVersionAnnotation: keyv1.VersionBundleVersion(customObject),
			},
			Labels: map[string]string{
				certificatetpr.ClusterIDLabel: keyv1.ClusterID(customObject),
				certificatetpr.ComponentLabel: keyv1.ClusterComponent(customObject),
			},
		},
		StringData: map[string]string{
			certificatetpr.CA.String():  "",
			certificatetpr.Crt.String(): "",
			certificatetpr.Key.String(): "",
		},
	}

	r.logger.Log("cluster", keyv1.ClusterID(customObject), "debug", "computed the desired secret")

	return secret, nil
}
