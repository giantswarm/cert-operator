package vaultcrt

import (
	"context"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultrole"

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

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "creating the secret in the Kubernetes API")

	// Issue certificates and fill the secret.
	{
		err := r.ensureVaultRole(customObject)
		if err != nil {
			return microerror.Mask(err)
		}

		ca, crt, key, err := r.issueCertificate(customObject)
		if err != nil {
			return microerror.Mask(err)
		}

		secretToCreate.StringData = map[string]string{
			certificatetpr.CA.String():  ca,
			certificatetpr.Crt.String(): crt,
			certificatetpr.Key.String(): key,
		}
	}

	_, err = r.k8sClient.CoreV1().Secrets(r.namespace).Create(secretToCreate)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Log("cluster", key.ClusterID(customObject), "debug", "created the secret in the Kubernetes API")

	return nil
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
		Organizations: []string{
			key.ClusterComponent(customObject),
		},
		TTL: key.CrtTTL(customObject),
	}
	result, err := r.vaultCrt.Create(c)
	if err != nil {
		return "", "", "", microerror.Mask(err)
	}

	return result.CA, result.Crt, result.Key, nil
}
