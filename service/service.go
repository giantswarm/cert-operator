// Package service implements business logic to issue certificates for clusters
// running on the Giantnetes platform.
package service

import (
	"fmt"
	"sync"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	k8sutil "github.com/giantswarm/operatorkit/client/k8s"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	vaultutil "github.com/giantswarm/cert-operator/client/vault"
	"github.com/giantswarm/cert-operator/flag"
	"github.com/giantswarm/cert-operator/service/ca"
	"github.com/giantswarm/cert-operator/service/crt"
	"github.com/giantswarm/cert-operator/service/version"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	KubernetesClient *kubernetes.Clientset
	Logger           micrologger.Logger
	VaultClient      *vaultapi.Client

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper

	Description string
	GitCommit   string
	Name        string
	Source      string
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		KubernetesClient: nil,
		Logger:           nil,
		VaultClient:      nil,

		// Settings.
		Flag:  nil,
		Viper: nil,

		Description: "",
		GitCommit:   "",
		Name:        "",
		Source:      "",
	}
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	config.Logger.Log("debug", fmt.Sprintf("creating cert-operator with config: %#v", config))

	var err error

	var k8sClient kubernetes.Interface
	{
		k8sConfig := k8sutil.Config{
			Logger: config.Logger,

			TLS: k8sutil.TLSClientConfig{
				CAFile:  config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile),
				CrtFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile),
				KeyFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile),
			},
			Address:   config.Viper.GetString(config.Flag.Service.Kubernetes.Address),
			InCluster: config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster),
		}

		k8sClient, err = k8sutil.NewClient(k8sConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var vaultClient *vaultapi.Client
	{
		vaultConfig := vaultutil.Config{
			Flag:  config.Flag,
			Viper: config.Viper,
		}

		vaultClient, err = vaultutil.NewClient(vaultConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var caService *ca.Service
	{
		caConfig := ca.DefaultConfig()
		caConfig.Flag = config.Flag
		caConfig.Logger = config.Logger
		caConfig.VaultClient = vaultClient
		caConfig.Viper = config.Viper

		caService, err = ca.New(caConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var crtService *crt.Service
	{
		crtConfig := crt.DefaultConfig()
		crtConfig.Flag = config.Flag
		crtConfig.CAService = caService
		crtConfig.K8sClient = k8sClient
		crtConfig.Logger = config.Logger
		crtConfig.VaultClient = vaultClient
		crtConfig.Viper = config.Viper

		crtService, err = crt.New(crtConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var versionService *version.Service
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Description = config.Description
		versionConfig.GitCommit = config.GitCommit
		versionConfig.Name = config.Name
		versionConfig.Source = config.Source

		versionService, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newService := &Service{
		// Dependencies.
		Crt:     crtService,
		Version: versionService,

		// Internals
		bootOnce: sync.Once{},
	}

	return newService, nil
}

type Service struct {
	// Dependencies.
	Crt     *crt.Service
	Version *version.Service

	// Internals.
	bootOnce sync.Once
}

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		s.Crt.Boot()
	})
}
