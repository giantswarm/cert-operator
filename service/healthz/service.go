package healthz

import (
	"context"
	"sync"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	vaultapi "github.com/hashicorp/vault/api"
)

// Config represents the configuration used to create a healthz service.
type Config struct {
	// Dependencies.
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client
}

// DefaultConfig provides a default configuration to create a new healthz
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:      nil,
		VaultClient: nil,
	}
}

// New creates a new configured healthz service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}
	if config.VaultClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "vault client must not be empty")
	}

	newService := &Service{
		// Dependencies.
		vaultClient: config.VaultClient,
		logger:      config.Logger,

		// Internals.
		bootOnce: sync.Once{},
	}

	return newService, nil
}

// Service implements the healthz service interface.
type Service struct {
	// Dependencies.
	logger      micrologger.Logger
	vaultClient *vaultapi.Client

	// Internals.
	bootOnce sync.Once
}

// Check implements the health check which lists the mounts for the sys backend.
// This checks that the operator can connect to the Vault API and the token is
// valid.
func (s *Service) Check(ctx context.Context, request Request) (*Response, error) {
	sysBackend := s.vaultClient.Sys()

	_, err := sysBackend.ListMounts()
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return DefaultResponse(), nil
}
