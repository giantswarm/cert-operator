package vaultcrt

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	apiv1 "k8s.io/api/core/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cert-operator/service/certconfig/v2/key"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "computing the desired secret")

	hash, err := key.CustomObjectHash(customObject)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// NOTE that the actual secret content here is left blank because only the
	// issuer backend, e.g. Vault, can generate certificates. This has to be
	// considered when computing the create, delete and update state.
	secret := &apiv1.Secret{
		ObjectMeta: apismetav1.ObjectMeta{
			Name: key.SecretName(customObject),
			Annotations: map[string]string{
				ConfigHashAnnotation:           hash,
				UpdateTimestampAnnotation:      r.currentTimeFactory().In(time.UTC).Format(UpdateTimestampLayout),
				VersionBundleVersionAnnotation: key.VersionBundleVersion(customObject),
			},
			Labels: map[string]string{
				key.ClusterIDLabel: key.ClusterID(customObject),
				key.ComponentLabel: key.ClusterComponent(customObject),
			},
		},
		StringData: map[string]string{
			key.CAID:  "",
			key.CrtID: "",
			key.KeyID: "",
		},
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "computed the desired secret")

	return secret, nil
}
