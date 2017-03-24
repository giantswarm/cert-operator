package k8s

import (
	"net/url"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/cert-operator/flag"
)

const (
	// Maximum QPS to the master from this client.
	MaxQPS = 100
	// Maximum burst for throttle.
	MaxBurst = 100
)

type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper
}

func newRawClientConfig(config Config) *rest.Config {
	tlsClientConfig := rest.TLSClientConfig{
		CertFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CertFile),
		KeyFile:  config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile),
		CAFile:   config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile),
	}
	rawClientConfig := &rest.Config{
		Host:            config.Viper.GetString(config.Flag.Service.Kubernetes.Address),
		QPS:             MaxQPS,
		Burst:           MaxBurst,
		TLSClientConfig: tlsClientConfig,
	}

	return rawClientConfig
}

func getRawClientConfig(config Config) (*rest.Config, error) {
	var rawClientConfig *rest.Config
	var err error

	address := config.Viper.GetString(config.Flag.Service.Kubernetes.Address)

	if config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster) {
		config.Logger.Log("debug", "creating in-cluster config")
		rawClientConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		if address != "" {
			config.Logger.Log("debug", "using explicit api server")
			rawClientConfig.Host = address
		}

	} else {
		if address == "" {
			return nil, microerror.MaskAnyf(invalidConfigError, "kubernetes address must not be empty")
		}

		config.Logger.Log("debug", "creating out-cluster config")

		// Kubernetes listen URL.
		_, err := url.Parse(address)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}

		rawClientConfig = newRawClientConfig(config)
	}

	return rawClientConfig, nil
}

func NewClient(config Config) (kubernetes.Interface, error) {
	rawClientConfig, err := getRawClientConfig(config)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(rawClientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
