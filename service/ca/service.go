package ca

import (
	"sync"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"

	"github.com/giantswarm/cert-operator/flag"
)

// Config represents the configuration used to create a CA service.
type Config struct {
	// Dependencies.
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper
}

// DefaultConfig provides a default configuration to create a new CA service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:      nil,
		VaultClient: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// New creates a new configured CA service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}
	if config.VaultClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "vault client must not be empty")
	}

	// Settings.
	if config.Flag == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "viper must not be empty")
	}

	newService := &Service{
		Config: config,

		// Internals
		bootOnce: sync.Once{},
	}

	return newService, nil
}

// Service implements the CA service interface.
type Service struct {
	Config

	// Internals.
	bootOnce sync.Once
}
