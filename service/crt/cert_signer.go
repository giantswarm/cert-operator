package crt

import (
	"strings"

	"github.com/giantswarm/certctl/service/cert-signer"
	"github.com/giantswarm/certctl/service/spec"
	"github.com/giantswarm/certificatetpr"
	microerror "github.com/giantswarm/microkit/error"
)

// Issue generates a certificate using the PKI backend signed by the certificate
// authority associated with the configured cluster ID. The certificate is saved
// as a set of k8s secrets.
func (s *Service) Issue(cert *certificatetpr.CustomObject) error {
	newCertSignerConfig := certsigner.DefaultConfig()
	newCertSignerConfig.VaultClient = s.Config.VaultClient

	newCertSigner, err := certsigner.New(newCertSignerConfig)
	if err != nil {
		return microerror.MaskAny(err)
	}

	// Generate a new signed certificate.
	newIssueConfig := spec.IssueConfig{
		ClusterID:  cert.Spec.ClusterID,
		CommonName: cert.Spec.CommonName,
		IPSANs:     strings.Join(cert.Spec.IPSANs, ","),
		AltNames:   strings.Join(cert.Spec.AltNames, ","),
		TTL:        cert.Spec.TTL,
	}
	newIssueResponse, err := newCertSigner.Issue(newIssueConfig)
	if err != nil {
		return microerror.MaskAny(err)
	}

	// Save the certificate as a k8s secret.
	secret := certificateSecret{
		ClusterComponent: cert.Spec.ClusterComponent,
		CommonName:       cert.Spec.CommonName,
		Namespace:        cert.ObjectMeta.Namespace,
		IssueResponse:    newIssueResponse,
	}
	if err := s.CreateCertificate(secret); err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}
