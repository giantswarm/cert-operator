package v2

import (
	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/resource/metricsresource"
	"github.com/giantswarm/operatorkit/framework/resource/retryresource"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultpki"
	"github.com/giantswarm/vaultrole"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/cert-operator/service/certconfig/v2/key"
	vaultcrtresource "github.com/giantswarm/cert-operator/service/certconfig/v2/resources/vaultcrt"
	vaultpkiresource "github.com/giantswarm/cert-operator/service/certconfig/v2/resources/vaultpki"
	vaultroleresource "github.com/giantswarm/cert-operator/service/certconfig/v2/resources/vaultrole"
)

const (
	ResourceRetries uint64 = 3
)

type ResourceSetConfig struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
	VaultCrt  vaultcrt.Interface
	VaultPKI  vaultpki.Interface
	VaultRole vaultrole.Interface

	ExpirationThreshold   time.Duration
	HandledVersionBundles []string
	// Name is the project name.
	Name      string
	Namespace string
}

func NewResourceSet(config ResourceSetConfig) (*framework.ResourceSet, error) {
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
	if len(config.HandledVersionBundles) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.HandledVersionBundles must not be empty")
	}
	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Name must not be empty")
	}
	if config.Namespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Namespace must not be empty")
	}

	var err error

	var vaultCrtResource framework.Resource
	{
		c := vaultcrtresource.DefaultConfig()

		c.CurrentTimeFactory = func() time.Time { return time.Now() }
		c.K8sClient = config.K8sClient
		c.Logger = config.Logger
		c.VaultCrt = config.VaultCrt

		c.ExpirationThreshold = config.ExpirationThreshold
		c.Namespace = config.Namespace

		vaultCrtResource, err = vaultcrtresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultPKIResource framework.Resource
	{
		c := vaultpkiresource.DefaultConfig()

		c.Logger = config.Logger
		c.VaultPKI = config.VaultPKI

		vaultPKIResource, err = vaultpkiresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultRoleResource framework.Resource
	{
		c := vaultroleresource.DefaultConfig()

		c.Logger = config.Logger
		c.VaultRole = config.VaultRole

		vaultRoleResource, err = vaultroleresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// We create the list of resources and wrap each resource around some common
	// resources like metrics and retry resources.
	//
	// NOTE that the retry resources wrap the underlying resources first. The
	// wrapped resources are then wrapped around the metrics resource. That way
	// the metrics also consider execution times and execution attempts including
	// retries.
	//
	// NOTE that the order of the namespace resource is important. We have to
	// start the namespace resource at first because all other resources are
	// created within this namespace.
	var resources []framework.Resource
	{
		resources = []framework.Resource{
			vaultPKIResource,
			vaultRoleResource,
			vaultCrtResource,
		}

		retryWrapConfig := retryresource.WrapConfig{}
		retryWrapConfig.BackOffFactory = func() backoff.BackOff { return backoff.WithMaxTries(backoff.NewExponentialBackOff(), ResourceRetries) }
		retryWrapConfig.Logger = config.Logger
		resources, err = retryresource.Wrap(resources, retryWrapConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		metricsWrapConfig := metricsresource.WrapConfig{}
		metricsWrapConfig.Name = config.Name
		resources, err = metricsresource.Wrap(resources, metricsWrapConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	handlesFunc := func(obj interface{}) bool {
		kvmConfig, err := key.ToCustomObject(obj)
		if err != nil {
			return false
		}
		versionBundleVersion := key.VersionBundleVersion(kvmConfig)

		for _, v := range config.HandledVersionBundles {
			if versionBundleVersion == v {
				return true
			}
		}

		return false
	}

	var resourceSet *framework.ResourceSet
	{
		c := framework.ResourceSetConfig{
			Handles:   handlesFunc,
			Logger:    config.Logger,
			Resources: resources,
		}

		resourceSet, err = framework.NewResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return resourceSet, nil
}
