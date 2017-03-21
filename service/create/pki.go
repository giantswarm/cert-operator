package create

import (
	"strings"

	"github.com/giantswarm/certctl/service/pki"
	"github.com/giantswarm/certctl/service/token"
)

// CheckPKIBackend checks if there is a valid PKI backend in Vault
// for the configured cluster ID.
func (s *Service) CheckPKIBackend(clusterID string) bool {
	pkiService, err := s.getPKIService()
	if err != nil {
		return false
	}

	tokenService, err := s.getTokenService()
	if err != nil {
		return false
	}

	// Check PKI config.
	mounted, err := pkiService.IsMounted(clusterID)
	if !mounted || err != nil {
		return false
	}
	caGenerated, err := pkiService.IsCAGenerated(clusterID)
	if !caGenerated || err != nil {
		return false
	}
	roleCreated, err := pkiService.IsRoleCreated(clusterID)
	if !roleCreated || err != nil {
		return false
	}

	// Check token config.
	policyCreated, err := tokenService.IsPolicyCreated(clusterID)
func (s *Service) checkPKIPolicy(clusterID string) bool {
	service, err := s.getTokenService()
	if err != nil {
		return false
	}

	// Check if there is a PKI policy.
	policyCreated, err := service.IsPolicyCreated(clusterID)
	if !policyCreated || err != nil {
		return false
	}

	// PKI policy is valid.
	return true
}

func (s *Service) createPKIPolicy(cert CertificateSpec) error {
	service, err := s.getTokenService()
	if err != nil {
		return microerror.MaskAny(err)
	}

	// Create PKI policy TODO Use latest certctl with CreateConfig object.
	return service.CreatePolicy(cert.ClusterID)
}

// Get the common name for the Cluster CA by removing the prefix
// from the certificate common name.
func (s *Service) getClusterCA(cert CertificateSpec) string {
	url := strings.Split(cert.CommonName, ".")
	prefix := url[0] + "."

	return strings.Replace(cert.CommonName, prefix, "", 1)
}

// Get the allowed domains which are the Cluster CA common name and the
// alt names specified for the certificate.
func (s *Service) getAllowedDomainsForCA(cert CertificateSpec) string {
	ca := s.getClusterCA(cert)

	domains := []string{ca}
	domains = append(domains, cert.AltNames...)

	return strings.Join(domains, ",")
}

func (s *Service) getPKIService() (pki.Service, error) {
	pkiConfig := pki.DefaultServiceConfig()
	pkiConfig.VaultClient = s.Config.VaultClient

	return pki.NewService(pkiConfig)
}

func (s *Service) getTokenService() (token.Service, error) {
	tokenConfig := token.DefaultServiceConfig()
	tokenConfig.VaultClient = s.Config.VaultClient

	return token.NewService(tokenConfig)
}
