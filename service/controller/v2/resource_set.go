package v2

import (
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/resource/metricsresource"
	"github.com/giantswarm/operatorkit/controller/resource/retryresource"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultpki"
	"github.com/giantswarm/vaultrole"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/cert-operator/service/controller/v2/key"
	vaultcrtresource "github.com/giantswarm/cert-operator/service/controller/v2/resources/vaultcrt"
	vaultpkiresource "github.com/giantswarm/cert-operator/service/controller/v2/resources/vaultpki"
	vaultroleresource "github.com/giantswarm/cert-operator/service/controller/v2/resources/vaultrole"
)

type ResourceSetConfig struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
	VaultCrt  vaultcrt.Interface
	VaultPKI  vaultpki.Interface
	VaultRole vaultrole.Interface

	ExpirationThreshold time.Duration
	Namespace           string
	ProjectName         string
}

func NewResourceSet(config ResourceSetConfig) (*controller.ResourceSet, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.VaultCrt == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultCrt must not be empty")
	}
	if config.VaultPKI == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultPKI must not be empty")
	}
	if config.VaultRole == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultRole must not be empty")
	}

	if config.ExpirationThreshold == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ExpirationThreshold must not be empty")
	}
	if config.Namespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Namespace must not be empty")
	}
	if config.ProjectName == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ProjectName must not be empty")
	}

	var err error

	var vaultCrtResource controller.Resource
	{
		c := vaultcrtresource.DefaultConfig()

		c.CurrentTimeFactory = func() time.Time { return time.Now() }
		c.K8sClient = config.K8sClient
		c.Logger = config.Logger
		c.VaultCrt = config.VaultCrt

		c.ExpirationThreshold = config.ExpirationThreshold
		c.Namespace = config.Namespace

		ops, err := vaultcrtresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		vaultCrtResource, err = toCRUDResource(config.Logger, ops)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultPKIResource controller.Resource
	{
		c := vaultpkiresource.DefaultConfig()

		c.Logger = config.Logger
		c.VaultPKI = config.VaultPKI

		ops, err := vaultpkiresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		vaultPKIResource, err = toCRUDResource(config.Logger, ops)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultRoleResource controller.Resource
	{
		c := vaultroleresource.DefaultConfig()

		c.Logger = config.Logger
		c.VaultRole = config.VaultRole

		ops, err := vaultroleresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		vaultRoleResource, err = toCRUDResource(config.Logger, ops)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []controller.Resource{
		vaultPKIResource,
		vaultRoleResource,
		vaultCrtResource,
	}

	{
		c := retryresource.WrapConfig{
			Logger: config.Logger,
		}

		resources, err = retryresource.Wrap(resources, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	{
		c := metricsresource.WrapConfig{
			Name: config.ProjectName,
		}

		resources, err = metricsresource.Wrap(resources, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	handlesFunc := func(obj interface{}) bool {
		customObject, err := key.ToCustomObject(obj)
		if err != nil {
			return false
		}

		if key.VersionBundleVersion(customObject) == VersionBundle().Version {
			return true
		}
		// TODO remove this hack with the next version bundle version or as soon as
		// all certconfigs obtain a real version bundle version.
		if key.VersionBundleVersion(customObject) == "" {
			return true
		}

		return false
	}

	var resourceSet *controller.ResourceSet
	{
		c := controller.ResourceSetConfig{
			Handles:   handlesFunc,
			Logger:    config.Logger,
			Resources: resources,
		}

		resourceSet, err = controller.NewResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return resourceSet, nil
}

func toCRUDResource(logger micrologger.Logger, ops controller.CRUDResourceOps) (*controller.CRUDResource, error) {
	c := controller.CRUDResourceConfig{
		Logger: logger,
		Ops:    ops,
	}

	r, err := controller.NewCRUDResource(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}
