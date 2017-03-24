package create

import (
	"github.com/giantswarm/certificatetpr"
	microerror "github.com/giantswarm/microkit/error"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
)

// CreateCertificate saves the certificate as a k8s secret.
func (s *Service) CreateCertificate(cert certificateSecret) error {
	var err error

	secret := &v1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name: cert.CommonName,
		},
		StringData: map[string]string{
			"crt": cert.IssueResponse.Certificate,
			"key": cert.IssueResponse.PrivateKey,
			"ca":  cert.IssueResponse.IssuingCA,
		},
	}

	// Create the secret which should be idempotent.
	_, err = s.Config.K8sClient.Core().Secrets(cert.Namespace).Create(secret)
	if errors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}

// DeleteCertificate deletes the k8s secret that stores the certificate.
func (s *Service) DeleteCertificate(cert *certificatetpr.CustomObject) error {
	namespace := cert.ObjectMeta.Namespace
	secretName := cert.Spec.CommonName

	// Delete the secret which should be idempotent.
	err := s.Config.K8sClient.Core().Secrets(namespace).Delete(secretName, &v1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}
