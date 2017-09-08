package crt

import (
	"strings"

	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/certctl/service/cert-signer"
	"github.com/giantswarm/certctl/service/spec"
	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
)

const (
	issueMaxElapsedTime = 30 * time.Second
)

// IssueAndWait generates a certificate using the PKI backend for the cluster.
// If an error occurs an exponential backoff is used. After the max elapsed time
// the error is returned to the caller.
func (s *Service) IssueAndWait(cert certificatetpr.Spec) error {
	initBackoff := backoff.NewExponentialBackOff()
	initBackoff.MaxElapsedTime = issueMaxElapsedTime

	operation := func() error {
		err := s.Issue(cert)
		if err != nil {
			s.Logger.Log("info", "failed to issue cert - retrying")
			return microerror.Mask(err)
		}

		return nil
	}

	return backoff.Retry(operation, initBackoff)
}

// Issue generates a certificate using the PKI backend signed by the certificate
// authority associated with the configured cluster ID. The certificate is saved
// as a set of k8s secrets.
func (s *Service) Issue(cert certificatetpr.Spec) error {
	newCertSignerConfig := certsigner.DefaultConfig()
	newCertSignerConfig.VaultClient = s.Config.VaultClient

	newCertSigner, err := certsigner.New(newCertSignerConfig)
	if err != nil {
		return microerror.Mask(err)
	}

	// Generate a new signed certificate.
	newIssueConfig := spec.IssueConfig{
		ClusterID:  cert.ClusterID,
		CommonName: cert.CommonName,
		IPSANs:     strings.Join(cert.IPSANs, ","),
		AltNames:   strings.Join(cert.AltNames, ","),
		TTL:        cert.TTL,
	}
	newIssueResponse, err := newCertSigner.Issue(newIssueConfig)
	if err != nil {
		return microerror.Mask(err)
	}

	// Save the certificate as a k8s secret.
	secret := certificateSecret{
		Certificate:   cert,
		IssueResponse: newIssueResponse,
	}
	if err := s.CreateCertificate(secret); err != nil {
		return microerror.Mask(err)
	}

	return nil
}
