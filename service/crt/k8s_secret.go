package crt

import (
	"fmt"

	"github.com/giantswarm/certificatetpr"
	microerror "github.com/giantswarm/microkit/error"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
)

// CreateCertificate saves the certificate as a k8s secret.
func (s *Service) CreateCertificate(secret certificateSecret) error {
	var err error

	k8sSecret := &v1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name: getSecretName(secret.Certificate),
			Labels: map[string]string{
				certificatetpr.ClusterIDLabel: secret.Certificate.ClusterID,
				certificatetpr.ComponentLabel: secret.Certificate.ClusterComponent,
			},
		},
		StringData: map[string]string{
			certificatetpr.Crt.String(): secret.IssueResponse.Certificate,
			certificatetpr.Key.String(): secret.IssueResponse.PrivateKey,
			certificatetpr.CA.String():  secret.IssueResponse.IssuingCA,
		},
	}

	// Create the secret which should be idempotent.
	_, err = s.Config.K8sClient.Core().Secrets(v1.NamespaceDefault).Create(k8sSecret)
	if errors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}

// DeleteCertificate deletes the k8s secret that stores the certificate.
func (s *Service) DeleteCertificate(cert certificatetpr.Spec) error {
	// Delete the secret which should be idempotent.
	err := s.Config.K8sClient.Core().Secrets(v1.NamespaceDefault).Delete(getSecretName(cert), &v1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}

func getSecretName(cert certificatetpr.Spec) string {
	return fmt.Sprintf("%s-%s", cert.ClusterID, cert.ClusterComponent)
}
