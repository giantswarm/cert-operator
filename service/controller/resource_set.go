package controller

import (
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/crud"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/wrapper/retryresource"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultpki"
	"github.com/giantswarm/vaultrole"
	vaultapi "github.com/hashicorp/vault/api"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/cert-operator/service/controller/resources/vaultaccess"
	vaultcrtresource "github.com/giantswarm/cert-operator/service/controller/resources/vaultcrt"
	vaultpkiresource "github.com/giantswarm/cert-operator/service/controller/resources/vaultpki"
	vaultroleresource "github.com/giantswarm/cert-operator/service/controller/resources/vaultrole"
)

type ResourceSetConfig struct {
	K8sClient   kubernetes.Interface
	CtrlClient  client.Client
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client
	VaultCrt    vaultcrt.Interface
	VaultPKI    vaultpki.Interface
	VaultRole   vaultrole.Interface

	ExpirationThreshold time.Duration
	Namespace           string
	ProjectName         string
}

func NewResourceSet(config ResourceSetConfig) ([]resource.Interface, error) {
	var err error

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
			CtrlClient:         config.CtrlClient,
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

	return resources, nil
}

func toCRUDResource(logger micrologger.Logger, v crud.Interface) (*crud.Resource, error) {
	c := crud.ResourceConfig{
		CRUD:   v,
		Logger: logger,
	}

	r, err := crud.NewResource(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}
