package create

import (
	"fmt"

	microerror "github.com/giantswarm/microkit/error"
)

type issueResponse struct {
	Certificate  string
	PrivateKey   string
	IssuingCA    string
	SerialNumber string
}

func (s *Service) Issue(config CertificateSpec) (issueResponse, error) {
	// Create a client for issuing a new signed certificate.
	logicalStore := s.Config.VaultClient.Logical()

	// Generate a certificate for the PKI backend signed by the certificate
	// authority associated with the configured cluster ID.
	data := map[string]interface{}{
		"ttl":         config.TTL,
		"common_name": config.CommonName,
		"ip_sans":     config.IPSANs,
		"alt_names":   config.AltNames,
	}

	secret, err := logicalStore.Write(s.SignedPath(config.ClusterID), data)
	if err != nil {
		return issueResponse{}, microerror.MaskAny(err)
	}

	// Collect the certificate data from the secret response.
	vCrt, ok := secret.Data["certificate"]
	if !ok {
		return issueResponse{}, microerror.MaskAnyf(keyPairNotFoundError, "certificate missing")
	}
	crt := vCrt.(string)

	vKey, ok := secret.Data["private_key"]
	if !ok {
		return issueResponse{}, microerror.MaskAnyf(keyPairNotFoundError, "private key missing")
	}
	key := vKey.(string)

	vCA, ok := secret.Data["issuing_ca"]
	if !ok {
		return issueResponse{}, microerror.MaskAnyf(keyPairNotFoundError, "issuing CA missing")
	}
	ca := vCA.(string)

	vSerial, ok := secret.Data["serial_number"]
	if !ok {
		return issueResponse{}, microerror.MaskAnyf(keyPairNotFoundError, "serial number missing")
	}
	serial := vSerial.(string)

	newIssueResponse := issueResponse{
		Certificate:  crt,
		PrivateKey:   key,
		IssuingCA:    ca,
		SerialNumber: serial,
	}

	return newIssueResponse, nil
}

func (s *Service) SignedPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/issue/role-%s", clusterID, clusterID)
}
