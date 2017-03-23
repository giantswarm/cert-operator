package create

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/giantswarm/certificatetpr"
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/tools/cache"

	k8sutil "github.com/giantswarm/cert-operator/client/k8s"
	"github.com/giantswarm/cert-operator/flag"
)

const (
	CertificateListAPIEndpoint  string = "/apis/giantswarm.io/v1/certificates"
	CertificateWatchAPIEndpoint string = "/apis/giantswarm.io/v1/watch/certificates"

	// Period for re-synchronizing the list of objects in k8s watcher. 0 means that re-sync will be
	// delayed as long as possible, until the watch will be closed or timed out.
	resyncPeriod time.Duration = 0
)

// Config represents the configuration used to create a create service.
type Config struct {
	// Dependencies.
	K8sClient   kubernetes.Interface
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper
}

// DefaultConfig provides a default configuration to create a new create service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		K8sClient:   nil,
		Logger:      nil,
		VaultClient: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// New creates a new configured version service.
func New(config Config) (*Service, error) {
	// Dependencies.
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

	newService := &Service{
		Config: config,

		// Internals
		bootOnce: sync.Once{},
	}

	return newService, nil
}

// Service implements the version service interface.
type Service struct {
	Config

	// Internals.
	bootOnce sync.Once
}

// Boot starts the service and implements the watch for the certificate TPR.
func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		if err := s.createTPR(); err != nil {
			panic(fmt.Sprintf("could not create cluster resource: %#v", err))
		}
		s.Config.Logger.Log("info", "successfully created third-party resource")

		_, certInformer := cache.NewInformer(
			s.newCertificateListWatch(),
			&certificatetpr.CustomObject{},
			resyncPeriod,
			cache.ResourceEventHandlerFuncs{
				AddFunc:    s.handleCertAdd,
				DeleteFunc: s.handleCertDelete,
				UpdateFunc: s.handleCertUpdate,
			},
		)

		s.Config.Logger.Log("info", "starting watch")

		// Certificate informer lifecycle can be interrupted by putting a value into a "stop channel".
		// We aren't currently using that functionality, so we are passing a nil here.
		certInformer.Run(nil)
	})
}

// handleCertAdd issues a certificate using Vault for the certificate TPR. A PKI backend is
// setup for the Cluster ID if it does not yet exist.
func (s *Service) handleCertAdd(obj interface{}) {
	cert := obj.(*certificatetpr.CustomObject)
	s.Config.Logger.Log("info", fmt.Sprintf("creating certificate '%s'", cert.Spec.CommonName))

	if err := s.setupPKIBackend(cert.Spec); err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not setup pki backend '%#v'", err))
		return
	}
	if err := s.setupPKIPolicy(cert.Spec); err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not setup pki backend '%#v'", err))
		return
	}

	_, err := s.Issue(cert.Spec)
	if err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not issue cert '%#v'", err))
		return
	}
	s.Config.Logger.Log("info", fmt.Sprintf("certificate issued %s", cert.Spec.CommonName))
}

// handleCertDelete is not yet implemented.
func (s *Service) handleCertDelete(obj interface{}) {
	cert := obj.(*certificatetpr.CustomObject)
	s.Config.Logger.Log("info", fmt.Sprintf("deleting certificate '%s' is not implemented yet", cert.Spec.CommonName))
}

// handleCertUpdate is not yet implemented.
func (s *Service) handleCertUpdate(old, cur interface{}) {
	cert := cur.(*certificatetpr.CustomObject)
	s.Config.Logger.Log("info", fmt.Sprintf("updating certificate '%s' is not implemented yet", cert.Spec.CommonName))
}

// newCertificateListWatch returns a configured list watch for the certificate TPR.
func (s *Service) newCertificateListWatch() *cache.ListWatch {
	client := s.Config.K8sClient.Core().RESTClient()

	listWatch := &cache.ListWatch{
		ListFunc: func(options api.ListOptions) (runtime.Object, error) {
			req := client.Get().AbsPath(CertificateListAPIEndpoint)
			b, err := req.DoRaw()
			if err != nil {
				return nil, err
			}

			var c certificatetpr.List
			if err := json.Unmarshal(b, &c); err != nil {
				return nil, err
			}

			return &c, nil
		},

		WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
			req := client.Get().AbsPath(CertificateWatchAPIEndpoint)
			stream, err := req.Stream()
			if err != nil {
				return nil, err
			}

			watcher := watch.NewStreamWatcher(&k8sutil.CertificateDecoder{
				Stream: stream,
			})

			return watcher, nil
		},
	}
	return listWatch
}
