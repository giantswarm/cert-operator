// Package service implements business logic to TODO
package service

import (
	"fmt"
	"sync"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	k8sutil "github.com/giantswarm/cert-operator/client/k8s"
	"github.com/giantswarm/cert-operator/flag"
	"github.com/giantswarm/cert-operator/service/create"
	"github.com/giantswarm/cert-operator/service/version"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	KubernetesClient *kubernetes.Clientset
	Logger           micrologger.Logger

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
			Flag:   config.Flag,
			Viper:  config.Viper,
		}

		k8sClient, err = k8sutil.NewClient(k8sConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var createService *create.Service
	{
		createConfig := create.DefaultConfig()
		createConfig.Flag = config.Flag
		createConfig.K8sClient = k8sClient
		createConfig.Logger = config.Logger
		createConfig.Viper = config.Viper

		createService, err = create.New(createConfig)
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
		Create:  createService,
		Version: versionService,

		// Internals
		bootOnce: sync.Once{},
	}

	return newService, nil
}

type Service struct {
	// Dependencies.
	Create  *create.Service
	Version *version.Service

	// Internals.
	bootOnce sync.Once
}

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		s.Create.Boot()
	})
}
