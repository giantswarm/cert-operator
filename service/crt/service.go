package crt

import (
	"fmt"
	"sync"

	"github.com/giantswarm/certctl/service/spec"
	"github.com/giantswarm/certificatetpr"
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/giantswarm/operatorkit/tpr"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/giantswarm/cert-operator/flag"
	"github.com/giantswarm/cert-operator/service/ca"
)

const (
	TPRVersion     = "v1"
	TPRDescription = "Managed certificates on Kubernetes clusters"
)

// Config represents the configuration used to create a Crt service.
type Config struct {
	// Dependencies.
	CAService   *ca.Service
	Logger      micrologger.Logger
	K8sClient   kubernetes.Interface
	VaultClient *vaultapi.Client

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper
}

// certificateSecret stores a cert issued by Vault that will be stored as a k8s secret.
type certificateSecret struct {
	Certificate   certificatetpr.Spec
	IssueResponse spec.IssueResponse
}

// DefaultConfig provides a default configuration to create a new create service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		CAService:   nil,
		K8sClient:   nil,
		Logger:      nil,
		VaultClient: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// Service implements the Crt service interface.
type Service struct {
	Config

	// Internals.
	bootOnce sync.Once
	tpr      *tpr.TPR
}

// New creates a new configured Crt service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.CAService == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "ca service must not be empty")
	}
	if config.K8sClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "kubernetes client must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}
	if config.VaultClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "vault client must not be empty")
	}

	// Settings.
	if config.Flag == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "viper must not be empty")
	}

	tprConfig := tpr.Config{
		K8sClient:   config.K8sClient,
		Logger:      config.Logger,
		Name:        certificatetpr.Name,
		Version:     TPRVersion,
		Description: TPRDescription,
	}
	tpr, err := tpr.New(tprConfig)
	if err != nil {
		return nil, microerror.MaskAnyf(err, "creating TPR for %#v", tprConfig)
	}

	newService := &Service{
		Config: config,

		// Internals
		bootOnce: sync.Once{},
		tpr:      tpr,
	}

	return newService, nil
}

// Boot starts the service and implements the watch for the certificate TPR.
func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		err := s.tpr.CreateAndWait()
		if tpr.IsAlreadyExists(err) {
			s.Logger.Log("debug", "third party resource already exists")
		} else if err != nil {
			s.Logger.Log("error", fmt.Sprintf("%#v", err))
			return
		}

		s.Logger.Log("debug", "starting list/watch")

		newResourceEventHandler := &cache.ResourceEventHandlerFuncs{
			AddFunc:    s.addFunc,
			DeleteFunc: s.deleteFunc,
		}
		newZeroObjectFactory := &tpr.ZeroObjectFactoryFuncs{
			NewObjectFunc:     func() runtime.Object { return &certificatetpr.CustomObject{} },
			NewObjectListFunc: func() runtime.Object { return &certificatetpr.List{} },
		}

		s.tpr.NewInformer(newResourceEventHandler, newZeroObjectFactory).Run(nil)
	})
}

// addFunc issues a certificate using Vault for the certificate TPR. A PKI backend is
// setup for the Cluster ID if it does not yet exist.
func (s *Service) addFunc(obj interface{}) {
	cert := *obj.(*certificatetpr.CustomObject)
	s.Config.Logger.Log("debug", fmt.Sprintf("creating certificate '%s'", cert.Spec.CommonName))

	if err := s.Config.CAService.SetupPKIAndWait(cert.Spec); err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not setup PKI '%#v'", err))
		return
	}
	if err := s.IssueAndWait(cert.Spec); err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not issue cert '%#v'", err))
		return
	}

	s.Config.Logger.Log("info", fmt.Sprintf("certificate '%s' issued", cert.Spec.CommonName))
}

// deleteFunc deletes the k8s secret containing the certificate.
func (s *Service) deleteFunc(obj interface{}) {
	cert := *obj.(*certificatetpr.CustomObject)
	s.Config.Logger.Log("debug", fmt.Sprintf("deleting certificate '%s'", cert.Spec.CommonName))

	if err := s.DeleteCertificateAndWait(cert.Spec); err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not delete certificate '%#v'", err))
		return
	}

	s.Config.Logger.Log("info", fmt.Sprintf("certificate '%s' deleted", cert.Spec.CommonName))
}
