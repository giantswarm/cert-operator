package controller

import (
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultpki"
	"github.com/giantswarm/vaultrole"
	vaultapi "github.com/hashicorp/vault/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v2 "github.com/giantswarm/cert-operator/service/controller/v2"
)

type CertConfig struct {
	K8sClient   k8sclient.Interface
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client

	CATTL               string
	CRDLabelSelector    string
	CommonNameFormat    string
	ExpirationThreshold time.Duration
	Namespace           string
	ProjectName         string
}

func (c CertConfig) newInformerListOptions() metav1.ListOptions {
	listOptions := metav1.ListOptions{
		LabelSelector: c.CRDLabelSelector,
	}

	return listOptions
}

type Cert struct {
	*controller.Controller
}

func NewCert(config CertConfig) (*Cert, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}

	var err error

	var vaultCrt vaultcrt.Interface
	{
		c := vaultcrt.DefaultConfig()

		c.Logger = config.Logger
		c.VaultClient = config.VaultClient

		vaultCrt, err = vaultcrt.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultPKI vaultpki.Interface
	{
		c := vaultpki.Config{
			Logger:      config.Logger,
			VaultClient: config.VaultClient,

			CATTL:            config.CATTL,
			CommonNameFormat: config.CommonNameFormat,
		}

		vaultPKI, err = vaultpki.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultRole vaultrole.Interface
	{
		c := vaultrole.DefaultConfig()

		c.Logger = config.Logger
		c.VaultClient = config.VaultClient

		c.CommonNameFormat = config.CommonNameFormat

		vaultRole, err = vaultrole.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var v2ResourceSet *controller.ResourceSet
	{
		c := v2.ResourceSetConfig{
			K8sClient:   config.K8sClient.K8sClient(),
			Logger:      config.Logger,
			VaultClient: config.VaultClient,
			VaultCrt:    vaultCrt,
			VaultPKI:    vaultPKI,
			VaultRole:   vaultRole,

			ExpirationThreshold: config.ExpirationThreshold,
			Namespace:           config.Namespace,
			ProjectName:         config.ProjectName,
		}

		v2ResourceSet, err = v2.NewResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			CRD:       v1alpha1.NewCertConfigCRD(),
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Name:      config.ProjectName,
			ResourceSets: []*controller.ResourceSet{
				v2ResourceSet,
			},
			NewRuntimeObjectFunc: func() runtime.Object {
				return new(v1alpha1.CertConfig)
			},
		}

		operatorkitController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	c := &Cert{
		Controller: operatorkitController,
	}

	return c, nil
}
