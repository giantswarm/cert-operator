package vaultcrt

import (
	"context"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultrole"
	apiv1 "k8s.io/client-go/pkg/api/v1"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	secretToCreate, err := toSecret(createChange)
	if err != nil {
		return microerror.Mask(err)
	}

	if secretToCreate != nil {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "creating the secret in the Kubernetes API")

		_, err := r.k8sClient.CoreV1().Secrets(r.namespace).Create(secretToCreate)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "created the secret in the Kubernetes API")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "the secret does not need to be created in the Kubernetes API")
	}

	return nil
}

func (r *Resource) newCreateChange(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	currentSecret, err := toSecret(currentState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	desiredSecret, err := toSecret(desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "finding out if the secret has to be created")

	var secretToCreate *apiv1.Secret
	if currentSecret == nil {
		secretToCreate = desiredSecret

		err := r.ensureVaultRole(customObject)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		ca, crt, key, err := r.issueCertificate(customObject)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		secretToCreate.StringData[certificatetpr.CA.String()] = ca
		secretToCreate.StringData[certificatetpr.Crt.String()] = crt
		secretToCreate.StringData[certificatetpr.Key.String()] = key
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "found out if the secret has to be created")

	return secretToCreate, nil
}

func (r *Resource) ensureVaultRole(customObject certificatetpr.CustomObject) error {
	c := vaultrole.ExistsConfig{
		ID: key.ClusterID(customObject),
		Organizations: []string{
			key.ClusterComponent(customObject),
		},
	}
	exists, err := r.vaultRole.Exists(c)
	if err != nil {
		return microerror.Mask(err)
	}

	if !exists {
		c := vaultrole.CreateConfig{
			AllowBareDomains: key.AllowBareDomains(customObject),
			AllowSubdomains:  AllowSubDomains,
			AltNames:         key.AltNames(customObject),
			ID:               key.ClusterID(customObject),
			Organizations: []string{
				key.ClusterComponent(customObject),
			},
			TTL: key.RoleTTL(customObject),
		}
		err := r.vaultRole.Create(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (r *Resource) issueCertificate(customObject certificatetpr.CustomObject) (string, string, string, error) {

	c := vaultcrt.CreateConfig{
		AltNames:   key.AltNames(customObject),
		CommonName: key.CommonName(customObject),
		ID:         key.ClusterID(customObject),
		IPSANs:     key.IPSANs(customObject),
		Organizations: append([]string{key.ClusterComponent(customObject)},
			key.Organizations(customObject)...),
		TTL: key.CrtTTL(customObject),
	}
	result, err := r.vaultCrt.Create(c)
	if err != nil {
		return "", "", "", microerror.Mask(err)
	}

	return result.CA, result.Crt, result.Key, nil
}
