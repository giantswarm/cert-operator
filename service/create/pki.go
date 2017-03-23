package create

import (
	"fmt"
	"strings"

	microerror "github.com/giantswarm/microkit/error"

	"github.com/giantswarm/certctl/service/pki"
	"github.com/giantswarm/certctl/service/token"
)

// setupPKIBackend creates a PKI backend if one does not exist for the cluster.
func (s *Service) setupPKIBackend(cert CertificateSpec) error {
	var err error

	service, err := s.getPKIService()
	if err != nil {
		return microerror.MaskAny(err)
	}

	isValid, err := service.VerifyPKISetup(cert.ClusterID)
	if err != nil {
		return microerror.MaskAny(err)
	}
	if isValid {
		s.Config.Logger.Log("debug", fmt.Sprintf("PKI backend already exists for cluster %s", cert.ClusterID))
		return nil
	}

	caCommonName := s.Config.getCACommonName(cert)

	// Create PKI backend
	config := pki.CreateConfig{
		ClusterID:        cert.ClusterID,
		CommonName:       caCommonName,
		AllowedDomains:   getAllowedDomainsForCA(caCommonName, cert),
		AllowBareDomains: cert.AllowBareDomains,
		TTL:              s.Config.Viper.GetString(s.Config.Flag.Vault.PKI.CATTL),
	}

	s.Config.Logger.Log("debug", fmt.Sprintf("Creating PKI backend for cluster %s", cert.ClusterID))
	return service.Create(config)
}

// setupPKIPolicy creates a PKI policy if one does not exist for the cluster.
func (s *Service) setupPKIPolicy(cert CertificateSpec) error {
	var err error

	service, err := s.getTokenService()
	if err != nil {
		return microerror.MaskAny(err)
	}

	isCreated, err := service.IsPolicyCreated(cert.ClusterID)
	if err != nil {
		return microerror.MaskAny(err)
	}
	if isCreated {
		s.Config.Logger.Log("debug", fmt.Sprintf("PKI policy already exists for cluster %s", cert.ClusterID))
		return nil
	}

	s.Config.Logger.Log("debug", fmt.Sprintf("Creating PKI policy for cluster %s", cert.ClusterID))
	return service.CreatePolicy(cert.ClusterID)
}

func (config Config) getCACommonName(cert CertificateSpec) string {
	commonNameFormat := config.Viper.GetString(config.Flag.Vault.PKI.CommonNameFormat)
	commonName := fmt.Sprintf(commonNameFormat, cert.ClusterID)

	return commonName
}

func getAllowedDomainsForCA(caCommonName string, cert CertificateSpec) string {
	domains := []string{caCommonName}
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
