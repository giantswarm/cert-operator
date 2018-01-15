package vaultrolev1

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/vaultrole"
)

const (
	Name = "vaultrolev1"
)

type Config struct {
	Logger    micrologger.Logger
	VaultRole vaultrole.Interface
}

func DefaultConfig() Config {
	return Config{
		Logger:    nil,
		VaultRole: nil,
	}
}

type Resource struct {
	logger    micrologger.Logger
	vaultRole vaultrole.Interface
}

func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.VaultRole == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultRole must not be empty")
	}

	r := &Resource{
		logger: config.Logger.With(
			"resource", Name,
		),
		vaultRole: config.VaultRole,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) Underlying() framework.Resource {
	return r
}
