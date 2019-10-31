package v2

import (
	"time"

	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/resource"
	"github.com/giantswarm/operatorkit/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/resource/wrapper/retryresource"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultpki"
	"github.com/giantswarm/vaultrole"
	vaultapi "github.com/hashicorp/vault/api"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/cert-operator/service/controller/v2/key"
	"github.com/giantswarm/cert-operator/service/controller/v2/resources/cleaning"
	"github.com/giantswarm/cert-operator/service/controller/v2/resources/vaultaccess"
	vaultcrtresource "github.com/giantswarm/cert-operator/service/controller/v2/resources/vaultcrt"
	vaultpkiresource "github.com/giantswarm/cert-operator/service/controller/v2/resources/vaultpki"
	vaultroleresource "github.com/giantswarm/cert-operator/service/controller/v2/resources/vaultrole"
)

type ResourceSetConfig struct {
	G8sClinet   versioned.Interface
	K8sClient   kubernetes.Interface
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client
	VaultCrt    vaultcrt.Interface
	VaultPKI    vaultpki.Interface
	VaultRole   vaultrole.Interface

	ExpirationThreshold time.Duration
	Namespace           string
	ProjectName         string
}

func NewResourceSet(config ResourceSetConfig) (*controller.ResourceSet, error) {
	var err error
	var cleaningResource resource.Interface
	{
		c := cleaning.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		cleaningResource, err = cleaning.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultAccessResource resource.Interface
	{
		c := vaultaccess.Config{
			Logger:      config.Logger,
			VaultClient: config.VaultClient,
		}

		vaultAccessResource, err = vaultaccess.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultCrtResource resource.Interface
	{
		c := vaultcrtresource.Config{
			CurrentTimeFactory: func() time.Time { return time.Now() },
			K8sClient:          config.K8sClient,
			Logger:             config.Logger,
			VaultCrt:           config.VaultCrt,

			ExpirationThreshold: config.ExpirationThreshold,
			Namespace:           config.Namespace,
		}

		ops, err := vaultcrtresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		vaultCrtResource, err = toCRUDResource(config.Logger, ops)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultPKIResource resource.Interface
	{
		c := vaultpkiresource.Config{
			Logger:   config.Logger,
			VaultPKI: config.VaultPKI,
		}

		ops, err := vaultpkiresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		vaultPKIResource, err = toCRUDResource(config.Logger, ops)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultRoleResource resource.Interface
	{
		c := vaultroleresource.Config{
			Logger:    config.Logger,
			VaultRole: config.VaultRole,
		}

		ops, err := vaultroleresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		vaultRoleResource, err = toCRUDResource(config.Logger, ops)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		cleaningResource,
		vaultAccessResource,
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
		c := metricsresource.WrapConfig{}

		resources, err = metricsresource.Wrap(resources, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	handlesFunc := func(obj interface{}) bool {
		cr, err := key.ToCustomObject(obj)
		if err != nil {
			return false
		}

		if key.VersionBundleVersion(cr) == VersionBundle().Version {
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
