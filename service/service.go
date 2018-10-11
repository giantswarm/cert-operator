// Package service implements business logic to create Kubernetes resources
// against the Kubernetes API.
package service

import (
	"context"
	"sync"

	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	clientvault "github.com/giantswarm/cert-operator/client/vault"
	"github.com/giantswarm/cert-operator/flag"
	"github.com/giantswarm/cert-operator/service/collector"
	"github.com/giantswarm/cert-operator/service/controller"
)

type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Description string
	GitCommit   string
	ProjectName string
	Source      string
}

func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,

		Description: "",
		GitCommit:   "",
		ProjectName: "",
		Source:      "",
	}
}

type Service struct {
	Version *version.Service

	bootOnce          sync.Once
	certController    *controller.Cert
	operatorCollector *collector.Collector
}

func New(config Config) (*Service, error) {
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Viper must not be empty")
	}

	var err error

	var restConfig *rest.Config
	{
		c := k8srestconfig.Config{
			Logger: config.Logger,

			Address:   config.Viper.GetString(config.Flag.Service.Kubernetes.Address),
			InCluster: config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster),
			TLS: k8srestconfig.TLSClientConfig{
				CAFile:  config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile),
				CrtFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile),
				KeyFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile),
			},
		}

		restConfig, err = k8srestconfig.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	g8sClient, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	k8sExtClient, err := apiextensionsclient.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var vaultClient *vaultapi.Client
	{
		vaultConfig := clientvault.Config{
			Flag:  config.Flag,
			Viper: config.Viper,
		}

		vaultClient, err = clientvault.NewClient(vaultConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var certController *controller.Cert
	{
		c := controller.CertConfig{
			G8sClient:    g8sClient,
			K8sClient:    k8sClient,
			K8sExtClient: k8sExtClient,
			Logger:       config.Logger,
			VaultClient:  vaultClient,

			CATTL:               config.Viper.GetString(config.Flag.Service.Vault.Config.PKI.CA.TTL),
			CRDLabelSelector:    config.Viper.GetString(config.Flag.Service.CRD.LabelSelector),
			CommonNameFormat:    config.Viper.GetString(config.Flag.Service.Vault.Config.PKI.CommonName.Format),
			ExpirationThreshold: config.Viper.GetDuration(config.Flag.Service.Resource.VaultCrt.ExpirationThreshold),
			Namespace:           config.Viper.GetString(config.Flag.Service.Resource.VaultCrt.Namespace),
			ProjectName:         config.ProjectName,
		}

		certController, err = controller.NewCert(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var operatorCollector *collector.Collector
	{
		c := collector.Config{
			Logger:      config.Logger,
			VaultClient: vaultClient,
		}

		operatorCollector, err = collector.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		c := version.Config{
			Description:    config.Description,
			GitCommit:      config.GitCommit,
			Name:           config.ProjectName,
			Source:         config.Source,
			VersionBundles: NewVersionBundles(),
		}

		versionService, err = version.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Service{
		Version: versionService,

		bootOnce:          sync.Once{},
		certController:    certController,
		operatorCollector: operatorCollector,
	}

	return s, nil
}

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		go s.certController.Boot()
		go s.operatorCollector.Boot(context.Background())
	})
}
