package crt

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	deleteSecretMaxElapsedTime = 30 * time.Second
)

// CreateCertificate saves the certificate as a k8s secret.
func (s *Service) CreateCertificate(secret certificateSecret) error {
	var err error

	k8sSecret := &v1.Secret{
		ObjectMeta: apismetav1.ObjectMeta{
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
	if apierrors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// DeleteCertificateAndWait tries to delete the k8s secret. If an error occurs
// an exponential backoff is used. After the max elapsed time the error will be
// returned to the caller. The secret deletion is idempotent so no error is
// returned if the secret has already been deleted.
func (s *Service) DeleteCertificateAndWait(cert certificatetpr.Spec) error {
	initBackoff := backoff.NewExponentialBackOff()
	initBackoff.MaxElapsedTime = deleteSecretMaxElapsedTime

	operation := func() error {
		err := s.DeleteCertificate(cert)
		if err != nil {
			s.Logger.Log("info", "failed to delete secret - retrying")
			return microerror.Mask(err)
		}

		return nil
	}

	return backoff.Retry(operation, initBackoff)
}

// DeleteCertificate deletes the k8s secret that stores the certificate. The secret
// deletion is idempotent so no error is returned if the secret has already
// been deleted.
func (s *Service) DeleteCertificate(cert certificatetpr.Spec) error {
	err := s.Config.K8sClient.Core().Secrets(v1.NamespaceDefault).Delete(getSecretName(cert), &apismetav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func getSecretName(cert certificatetpr.Spec) string {
	return fmt.Sprintf("%s-%s", cert.ClusterID, cert.ClusterComponent)
}
