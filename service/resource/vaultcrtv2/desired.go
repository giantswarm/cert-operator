package vaultcrtv2

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	apiv1 "k8s.io/api/core/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/cert-operator/service/keyv2"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	customObject, err := keyv2.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "computing the desired secret")

	hash, err := keyv2.CustomObjectHash(customObject)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// NOTE that the actual secret content here is left blank because only the
	// issuer backend, e.g. Vault, can generate certificates. This has to be
	// considered when computing the create, delete and update state.
	secret := &apiv1.Secret{
		ObjectMeta: apismetav1.ObjectMeta{
			Name: keyv2.SecretName(customObject),
			Annotations: map[string]string{
				ConfigHashAnnotation:           hash,
				UpdateTimestampAnnotation:      r.currentTimeFactory().In(time.UTC).Format(UpdateTimestampLayout),
				VersionBundleVersionAnnotation: keyv2.VersionBundleVersion(customObject),
			},
			Labels: map[string]string{
				keyv2.ClusterIDLabel: keyv2.ClusterID(customObject),
				keyv2.ComponentLabel: keyv2.ClusterComponent(customObject),
			},
		},
		StringData: map[string]string{
			keyv2.CAID:  "",
			keyv2.CrtID: "",
			keyv2.KeyID: "",
		},
	}

	r.logger.Log("cluster", keyv2.ClusterID(customObject), "debug", "computed the desired secret")

	return secret, nil
}
