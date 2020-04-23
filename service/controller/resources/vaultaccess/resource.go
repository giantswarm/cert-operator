package vaultaccess

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	vaultapi "github.com/hashicorp/vault/api"
)

const (
	Name = "vaultaccess"
)

type Config struct {
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client
}

type Resource struct {
	logger      micrologger.Logger
	vaultClient *vaultapi.Client
}

func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.VaultClient must not be empty", config)
	}

	r := &Resource{
		logger:      config.Logger,
		vaultClient: config.VaultClient,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
