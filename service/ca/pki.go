package ca

import (
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/certctl/service/pki"
	"github.com/giantswarm/certctl/service/token"
	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
)

const (
	setupPKIMaxElapsedTime = 30 * time.Second
)

// SetupPKIAndWait creates a PKI backend and policy if one does not exist for
// the cluster. If an error occurs an exponential backoff is used. After the
// max elapsed time the error is returned to the caller.
func (s *Service) SetupPKIAndWait(cert certificatetpr.Spec) error {
	initBackoff := backoff.NewExponentialBackOff()
	initBackoff.MaxElapsedTime = setupPKIMaxElapsedTime

	operation := func() error {
		err := s.SetupPKI(cert)
		if err != nil {
			s.Logger.Log("info", "failed to setup PKI - retrying")
			return microerror.Mask(err)
		}

		return nil
	}

	return backoff.Retry(operation, initBackoff)
}

// SetupPKI creates a PKI backend and policy if one does not exist for the cluster.
func (s *Service) SetupPKI(cert certificatetpr.Spec) error {
	s.Config.Logger.Log("debug", fmt.Sprintf("setting up PKI for cluster '%s'", cert.ClusterID))

	if err := s.setupPKIBackend(cert); err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not setup pki backend '%#v'", err))
		return err
	}
	if err := s.setupPKIPolicy(cert); err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not setup pki policy '%#v'", err))
		return err
	}

	s.Config.Logger.Log("debug", fmt.Sprintf("valid PKI exists for cluster '%s'", cert.ClusterID))
	return nil
}

// setupPKIBackend creates a PKI backend if one does not exist for the cluster.
func (s *Service) setupPKIBackend(cert certificatetpr.Spec) error {
	var err error

	service, err := s.getPKIService()
	if err != nil {
		return microerror.Mask(err)
	}

	isValid, err := service.VerifyPKISetup(cert.ClusterID)
	if err != nil {
		return microerror.Mask(err)
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
		TTL:              s.Config.Viper.GetString(s.Config.Flag.Service.Vault.Config.PKI.CA.TTL),
	}

	s.Config.Logger.Log("debug", fmt.Sprintf("Creating PKI backend for cluster %s", cert.ClusterID))
	return service.Create(config)
}

// setupPKIPolicy creates a PKI policy if one does not exist for the cluster.
func (s *Service) setupPKIPolicy(cert certificatetpr.Spec) error {
	var err error

	service, err := s.getTokenService()
	if err != nil {
		return microerror.Mask(err)
	}

	isCreated, err := service.IsPolicyCreated(cert.ClusterID)
	if err != nil {
		return microerror.Mask(err)
	}
	if isCreated {
		s.Config.Logger.Log("debug", fmt.Sprintf("PKI policy already exists for cluster %s", cert.ClusterID))
		return nil
	}

	s.Config.Logger.Log("debug", fmt.Sprintf("Creating PKI policy for cluster %s", cert.ClusterID))
	return service.CreatePolicy(cert.ClusterID)
}

func (config Config) getCACommonName(cert certificatetpr.Spec) string {
	commonNameFormat := config.Viper.GetString(config.Flag.Service.Vault.Config.PKI.CommonName.Format)
	commonName := fmt.Sprintf(commonNameFormat, cert.ClusterID)

	return commonName
}

func getAllowedDomainsForCA(caCommonName string, cert certificatetpr.Spec) string {
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
