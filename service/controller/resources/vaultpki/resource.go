package vaultpki

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/vaultpki"
)

const (
	Name = "vaultpki"
)

type Config struct {
	Logger   micrologger.Logger
	VaultPKI vaultpki.Interface
}

type Resource struct {
	logger   micrologger.Logger
	vaultPKI vaultpki.Interface
}

func New(config Config) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.VaultPKI == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.VaultPKI must not be empty", config)
	}

	r := &Resource{
		logger:   config.Logger,
		vaultPKI: config.VaultPKI,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
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
