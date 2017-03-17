package create

import (
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
	if !policyCreated || err != nil {
		return false
	}

	// PKI config is valid.
	return true
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
