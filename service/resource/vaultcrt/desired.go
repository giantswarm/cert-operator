package vaultcrt

import (
	"context"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}, deleted bool) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "computing the desired secret")

	var secret *apiv1.Secret
	if !deleted {
		secret = &apiv1.Secret{
			ObjectMeta: apismetav1.ObjectMeta{
				Name: key.SecretName(customObject),
				Labels: map[string]string{
					certificatetpr.ClusterIDLabel: key.ClusterID(customObject),
					certificatetpr.ComponentLabel: key.ClusterComponent(customObject),
				},
			},
		}
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "computed the desired secret")

	return secret, nil
}
