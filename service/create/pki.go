package create

import (
	"fmt"
	"strings"

	microerror "github.com/giantswarm/microkit/error"

	"github.com/giantswarm/certctl/service/pki"
	"github.com/giantswarm/certctl/service/token"
)


func (s *Service) checkPKIBackend(clusterID string) bool {
	service, err := s.getPKIService()
	if err != nil {
		return false
	}

	// Check PKI config.
	mounted, err := service.IsMounted(clusterID)
	if !mounted || err != nil {
		return false
	}
	caGenerated, err := service.IsCAGenerated(clusterID)
	if !caGenerated || err != nil {
		return false
	}
	roleCreated, err := service.IsRoleCreated(clusterID)
	if !roleCreated || err != nil {
		return false
	}

	// PKI config is valid.
	return true
}

func (s *Service) createPKIBackend(cert CertificateSpec) error {
	var err error

	service, err := s.getPKIService()
	if err != nil {
		return microerror.MaskAny(err)
	}

	// Create PKI backend
	config := pki.CreateConfig{
		ClusterID:        cert.ClusterID,
		CommonName:       s.getCACommonName(cert),
		AllowedDomains:   s.getAllowedDomainsForCA(cert),
		AllowBareDomains: cert.AllowBareDomains,
		TTL:              s.Config.Viper.GetString(s.Config.Flag.Vault.PKI.CATTL),
	}

	return service.Create(config)
}

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

// Get the Common Name for the Cluster CA.
func (s *Service) getCACommonName(cert CertificateSpec) string {
	commonNameFormat := s.Config.Viper.GetString(s.Config.Flag.Vault.PKI.CommonNameFormat)
	commonName := fmt.Sprintf(commonNameFormat, cert.ClusterID)

	return commonName
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
