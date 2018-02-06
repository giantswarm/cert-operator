package vaultpki

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/vaultpki"
)

const (
	Name = "vaultpkiv2"
)

type Config struct {
	Logger   micrologger.Logger
	VaultPKI vaultpki.Interface
}

func DefaultConfig() Config {
	return Config{
		Logger:   nil,
		VaultPKI: nil,
	}
}

type Resource struct {
	logger   micrologger.Logger
	vaultPKI vaultpki.Interface
}

func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.VaultPKI == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultPKI must not be empty")
	}

	r := &Resource{
		logger: config.Logger.With(
			"resource", Name,
		),
		vaultPKI: config.VaultPKI,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) Underlying() framework.Resource {
	return r
}

func toVaultPKIState(v interface{}) (VaultPKIState, error) {
	if v == nil {
		return VaultPKIState{}, nil
	}

	vaultPKIState, ok := v.(VaultPKIState)
	if !ok {
		return VaultPKIState{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", VaultPKIState{}, v)
	}

	return vaultPKIState, nil
}
