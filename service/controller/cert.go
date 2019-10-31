package controller

import (
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8scrdclient"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/informer"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultpki"
	"github.com/giantswarm/vaultrole"
	vaultapi "github.com/hashicorp/vault/api"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/cert-operator/service/controller/v2"
)

type CertConfig struct {
	G8sClient    versioned.Interface
	K8sClient    kubernetes.Interface
	K8sExtClient apiextensionsclient.Interface
	Logger       micrologger.Logger
	VaultClient  *vaultapi.Client

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
	if config.G8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.G8sClient must not be empty")
	}

	var err error

	var crdClient *k8scrdclient.CRDClient
	{
		c := k8scrdclient.Config{
			K8sExtClient: config.K8sExtClient,
			Logger:       config.Logger,
		}

		crdClient, err = k8scrdclient.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

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

	var newInformer *informer.Informer
	{
		c := informer.Config{
			Logger:  config.Logger,
			Watcher: config.G8sClient.CoreV1alpha1().CertConfigs(""),

			ListOptions:  config.newInformerListOptions(),
			RateWait:     informer.DefaultRateWait,
			ResyncPeriod: informer.DefaultResyncPeriod,
		}

		newInformer, err = informer.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var v2ResourceSet *controller.ResourceSet
	{
		c := v2.ResourceSetConfig{
			G8sClinet:   config.G8sClient,
			K8sClient:   config.K8sClient,
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
			CRDClient: crdClient,
			Informer:  newInformer,
			Logger:    config.Logger,
			Name:      config.ProjectName,
			ResourceSets: []*controller.ResourceSet{
				v2ResourceSet,
			},
			RESTClient: config.G8sClient.CoreV1alpha1().RESTClient(),
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
