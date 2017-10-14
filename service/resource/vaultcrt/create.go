package vaultcrt

import (
	"context"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultrole"

	"github.com/giantswarm/cert-operator/service/key"
)

func (r *Resource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
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

func (r *Resource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	customObject, err := key.ToCustomObject(obj)
	if err != nil {
		return microerror.Mask(err)
	}
	secretToCreate, err := toSecret(createState)
	if err != nil {
		return microerror.Mask(err)
	}

	if secretToCreate != nil {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "creating the secret in the Kubernetes API")

		err := r.k8sClient.Core().Secrets(r.namespace).Create(secretToCreate)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "created the secret in the Kubernetes API")
	} else {
		r.logger.Log("cluster", key.ClusterID(customObject), "debug", "the secret does not need to be created in the Kubernetes API")
	}

	return nil
}

func (r *Resource) ensureVaultRole(customObject certificatetpr.CustomObject) error {
	// NOTE we do not set organizations yet because the TPR does not support it.
	c := vaultrole.ExistsConfig{
		ID:            key.ClusterID(customObject),
		Organizations: "",
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
			Organizations:    "",
			TTL:              key.RoleTTL(customObject),
		}
		err := r.vaultRole.Create(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return microerror.Mask(err)
}

func (r *Resource) issueCertificate(customObject certificatetpr.CustomObject) (string, string, string, error) {
	c := vaultcrt.CreateConfig{
		AltNames: key.AltNames(customObject),
		ID:       key.ClusterID(customObject),
		IPSANs:   key.IPSANs(customObject),
		TTL:      key.CrtTTL(customObject),
	}
	result, err := r.vaultCrt.Create(c)
	if err != nil {
		return microerror.Mask(err)
	}

	return result.CA, result.Crt, result.Key, nil
}
