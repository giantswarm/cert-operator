// Package service implements business logic to create Kubernetes resources
// against the Kubernetes API.
package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/k8sclient/k8srestconfig"
	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/vaultpki"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Version     string
}

type Service struct {
	Version *version.Service

	bootOnce          sync.Once
	certController    *controller.Cert
	operatorCollector *collector.Set
	logger            micrologger.Logger
	vaultPKI          vaultpki.Interface
	k8sClient         *k8sclient.Clients
}

func New(config Config) (*Service, error) {
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flag must not be empty", config)
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Viper must not be empty", config)
	}

	var err error

	var restConfig *rest.Config
	{
		c := k8srestconfig.Config{
			Logger: config.Logger,

			Address:    config.Viper.GetString(config.Flag.Service.Kubernetes.Address),
			InCluster:  config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster),
			KubeConfig: config.Viper.GetString(config.Flag.Service.Kubernetes.KubeConfig),
			TLS: k8srestconfig.ConfigTLS{
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

	var k8sClient *k8sclient.Clients
	{
		c := k8sclient.ClientsConfig{
			SchemeBuilder: k8sclient.SchemeBuilder{
				v1alpha1.AddToScheme,
			},
			Logger: config.Logger,

			RestConfig: restConfig,
		}

		k8sClient, err = k8sclient.NewClients(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
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
			K8sClient:   k8sClient,
			Logger:      config.Logger,
			VaultClient: vaultClient,

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

	var operatorCollector *collector.Set
	{
		c := collector.SetConfig{
			Logger:      config.Logger,
			VaultClient: vaultClient,
		}

		operatorCollector, err = collector.NewSet(c)
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
			Version:        config.Version,
			VersionBundles: NewVersionBundles(),
		}

		versionService, err = version.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	//TODO LH temporary to see things working
	var vaultPKI vaultpki.Interface
	{
		c := vaultpki.Config{
			Logger:      config.Logger,
			VaultClient: vaultClient,
			CATTL:       config.Viper.GetString(config.Flag.Service.Vault.Config.PKI.CA.TTL),

			CommonNameFormat: config.Viper.GetString(config.Flag.Service.Vault.Config.PKI.CommonName.Format),
		}

		vaultPKI, err = vaultpki.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Service{
		Version: versionService,

		bootOnce:          sync.Once{},
		certController:    certController,
		operatorCollector: operatorCollector,
		logger:            config.Logger,
		vaultPKI:          vaultPKI,
		k8sClient:         k8sClient,
	}

	return s, nil
}

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		go s.certController.Boot(context.Background())
		go s.operatorCollector.Boot(context.Background())
		s.CleanVault()
	})
}

func (s *Service) CleanVault() {
	//do some vault cleanup..
	//TODO LH refactor this somewhere more consistent
	go func() {
		for range time.Tick(24 * time.Second) {
			s.logger.Log("level", "debug", "message", "service.CleanVault() - started")
			backends, err := s.vaultPKI.GetBackends()
			if err != nil {
				microerror.Mask(err)
			}
			s.logger.Log("level", "debug", "message", fmt.Sprintf("backends : %+v", backends))
			s.logger.Log("level", "debug", "message", "service.CleanVault() - lookup clusters")
			clusters, err := s.k8sClient.G8sClient().InfrastructureV1alpha2().AWSClusters("default").List(metav1.ListOptions{})
			if err != nil {
				microerror.Mask(err)
			}
			s.logger.Log("level", "debug", "message", fmt.Sprintf("clusters : %+v", clusters))
		}
	}()
}
