package create

import (
	"strings"

	"github.com/giantswarm/certctl/service/cert-signer"
	"github.com/giantswarm/certctl/service/spec"
	microerror "github.com/giantswarm/microkit/error"
)

// Issue generates a certificate using the PKI backend signed by the certificate
// authority associated with the configured cluster ID.
func (s *Service) Issue(config CertificateSpec) (spec.IssueResponse, error) {
	var issueResp spec.IssueResponse

	// Create a certificate signer to generate a new signed certificate.
	newCertSignerConfig := certsigner.DefaultConfig()
	newCertSignerConfig.VaultClient = s.Config.VaultClient

	newCertSigner, err := certsigner.New(newCertSignerConfig)
	if err != nil {
		return issueResp, microerror.MaskAny(err)
	}

	// Generate a new signed certificate.
	newIssueConfig := spec.IssueConfig{
		ClusterID:  config.ClusterID,
		CommonName: config.CommonName,
		IPSANs:     strings.Join(config.IPSANs, ","),
		AltNames:   strings.Join(config.AltNames, ","),
		TTL:        config.TTL,
	}

	newIssueResponse, err := newCertSigner.Issue(newIssueConfig)
	if err != nil {
		return issueResp, microerror.MaskAny(err)
	}

	return newIssueResponse, err
}
