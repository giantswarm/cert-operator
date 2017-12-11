package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/operatorkit/client/k8sclient"
	"github.com/giantswarm/operatorkit/client/k8scrdclient"
	"github.com/giantswarm/operatorkit/client/k8sextclient"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/resource/metricsresource"
	"github.com/giantswarm/operatorkit/framework/resource/retryresource"
	"github.com/giantswarm/operatorkit/informer"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultpki"
	"github.com/giantswarm/vaultrole"
	vaultapi "github.com/hashicorp/vault/api"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	vaultutil "github.com/giantswarm/cert-operator/client/vault"
	"github.com/giantswarm/cert-operator/service/resource/vaultcrtv2"
	"github.com/giantswarm/cert-operator/service/resource/vaultpkiv2"
)

const (
	ResourceRetries uint64 = 3
)

const (
	CertConfigCleanupFinalizer = "cert-operator.giantswarm.io/custom-object-cleanup"
)

func newCRDFramework(config Config) (*framework.Framework, error) {
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Viper must not be empty")
	}

	var err error

	var k8sClient kubernetes.Interface
	{
		k8sConfig := k8sclient.DefaultConfig()

		k8sConfig.Logger = config.Logger

		k8sConfig.Address = config.Viper.GetString(config.Flag.Service.Kubernetes.Address)
		k8sConfig.InCluster = config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster)
		k8sConfig.TLS.CAFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile)
		k8sConfig.TLS.CrtFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile)
		k8sConfig.TLS.KeyFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile)

		k8sClient, err = k8sclient.New(k8sConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var k8sExtClient apiextensionsclient.Interface
	{
		c := k8sextclient.DefaultConfig()

		c.Logger = config.Logger

		c.Address = config.Viper.GetString(config.Flag.Service.Kubernetes.Address)
		c.InCluster = config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster)
		c.TLS.CAFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile)
		c.TLS.CrtFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile)
		c.TLS.KeyFile = config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile)

		k8sExtClient, err = k8sextclient.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var crdClient *k8scrdclient.CRDClient
	{
		c := k8scrdclient.DefaultConfig()

		c.K8sExtClient = k8sExtClient
		c.Logger = microloggertest.New()

		crdClient, err = k8scrdclient.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
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
			return nil, microerror.Mask(err)
		}
	}

	var vaultCrt vaultcrt.Interface
	{
		c := vaultcrt.DefaultConfig()

		c.Logger = config.Logger
		c.VaultClient = vaultClient

		vaultCrt, err = vaultcrt.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultPKI vaultpki.Interface
	{
		c := vaultpki.DefaultConfig()

		c.Logger = config.Logger
		c.VaultClient = vaultClient

		c.CATTL = config.Viper.GetString(config.Flag.Service.Vault.Config.PKI.CA.TTL)
		c.CommonNameFormat = config.Viper.GetString(config.Flag.Service.Vault.Config.PKI.CommonName.Format)

		vaultPKI, err = vaultpki.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultRole vaultrole.Interface
	{
		c := vaultrole.DefaultConfig()

		c.Logger = config.Logger
		c.VaultClient = vaultClient

		c.CommonNameFormat = config.Viper.GetString(config.Flag.Service.Vault.Config.PKI.CommonName.Format)

		vaultRole, err = vaultrole.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultCrtResource framework.Resource
	{
		c := vaultcrtv2.DefaultConfig()

		c.CurrentTimeFactory = func() time.Time { return time.Now() }
		c.K8sClient = k8sClient
		c.Logger = config.Logger
		c.VaultCrt = vaultCrt
		c.VaultRole = vaultRole

		c.ExpirationThreshold = config.Viper.GetDuration(config.Flag.Service.Resource.VaultCrt.ExpirationThreshold)
		c.Namespace = config.Viper.GetString(config.Flag.Service.Resource.VaultCrt.Namespace)

		vaultCrtResource, err = vaultcrtv2.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vaultPKIResource framework.Resource
	{
		c := vaultpkiv2.DefaultConfig()

		c.Logger = config.Logger
		c.VaultPKI = vaultPKI

		vaultPKIResource, err = vaultpkiv2.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// We create the list of resources and wrap each resource around some common
	// resources like metrics and retry resources.
	//
	// NOTE that the retry resources wrap the underlying resources first. The
	// wrapped resources are then wrapped around the metrics resource. That way
	// the metrics also consider execution times and execution attempts including
	// retries.
	//
	// NOTE that the order of the namespace resource is important. We have to
	// start the namespace resource at first because all other resources are
	// created within this namespace.
	var resources []framework.Resource
	{
		resources = []framework.Resource{
			vaultPKIResource,
			vaultCrtResource,
		}

		retryWrapConfig := retryresource.DefaultWrapConfig()
		retryWrapConfig.BackOffFactory = func() backoff.BackOff { return backoff.WithMaxTries(backoff.NewExponentialBackOff(), ResourceRetries) }
		retryWrapConfig.Logger = config.Logger
		resources, err = retryresource.Wrap(resources, retryWrapConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		metricsWrapConfig := metricsresource.DefaultWrapConfig()
		metricsWrapConfig.Name = config.Name
		resources, err = metricsresource.Wrap(resources, metricsWrapConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var clientSet *versioned.Clientset
	{
		var c *rest.Config

		if config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster) {
			config.Logger.Log("debug", "creating in-cluster config")

			c, err = rest.InClusterConfig()
			if err != nil {
				return nil, microerror.Mask(err)
			}
		} else {
			config.Logger.Log("debug", "creating out-cluster config")

			c = &rest.Config{
				Host: config.Viper.GetString(config.Flag.Service.Kubernetes.Address),
				TLSClientConfig: rest.TLSClientConfig{
					CAFile:   config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile),
					CertFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile),
					KeyFile:  config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile),
				},
			}
		}

		clientSet, err = versioned.NewForConfig(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// TODO remove after migration.
	go migrateTPRsToCRDs(config.Logger, clientSet)

	var newWatcherFactory informer.WatcherFactory
	{
		newWatcherFactory = func() (watch.Interface, error) {
			watcher, err := clientSet.CoreV1alpha1().CertConfigs("").Watch(apismetav1.ListOptions{})
			if err != nil {
				return nil, microerror.Mask(err)
			}

			return watcher, nil
		}
	}

	var newInformer *informer.Informer
	{
		informerConfig := informer.DefaultConfig()

		informerConfig.WatcherFactory = newWatcherFactory

		newInformer, err = informer.New(informerConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var crdFramework *framework.Framework
	{
		c := framework.DefaultConfig()

		c.CRD = v1alpha1.NewCertConfigCRD()
		c.CRDClient = crdClient
		c.Informer = newInformer
		c.Logger = config.Logger
		c.ResourceRouter = framework.DefaultResourceRouter(resources)

		crdFramework, err = framework.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return crdFramework, nil
}

func migrateTPRsToCRDs(logger micrologger.Logger, clientSet *versioned.Clientset) {
	logger.Log("debug", "start TPR migration")

	var err error

	// List all TPOs.
	var b []byte
	{
		e := "apis/giantswarm.io/v1/namespaces/default/certificates"
		b, err = clientSet.Discovery().RESTClient().Get().AbsPath(e).DoRaw()
		if err != nil {
			logger.Log("error", fmt.Sprintf("%#v", err))
			return
		}

		fmt.Printf("\n")
		fmt.Printf("b start\n")
		fmt.Printf("%s\n", b)
		fmt.Printf("b end\n")
		fmt.Printf("\n")
	}

	// Convert bytes into structure.
	var v *certificatetpr.List
	{
		v = &certificatetpr.List{}
		if err := json.Unmarshal(b, v); err != nil {
			logger.Log("error", fmt.Sprintf("%#v", err))
			return
		}

		fmt.Printf("\n")
		fmt.Printf("v start\n")
		fmt.Printf("%#v\n", v)
		fmt.Printf("v end\n")
		fmt.Printf("\n")
	}

	// Iterate over all TPOs.
	for _, tpo := range v.Items {
		// Compute CRO using TPO.
		var cro *v1alpha1.CertConfig
		{
			cro = &v1alpha1.CertConfig{}

			cro.TypeMeta.APIVersion = "core.giantswarm.io"
			cro.TypeMeta.Kind = "CertConfig"
			cro.ObjectMeta.Name = tpo.Name
			//cro.ObjectMeta.Finalizers = []string{
			//	CertConfigCleanupFinalizer,
			//}
			cro.Spec.Cert.AllowBareDomains = tpo.Spec.AllowBareDomains
			cro.Spec.Cert.AltNames = tpo.Spec.AltNames
			cro.Spec.Cert.ClusterComponent = tpo.Spec.ClusterComponent
			cro.Spec.Cert.ClusterID = tpo.Spec.ClusterID
			cro.Spec.Cert.CommonName = tpo.Spec.CommonName
			cro.Spec.Cert.IPSANs = tpo.Spec.IPSANs
			cro.Spec.Cert.Organizations = tpo.Spec.Organizations
			cro.Spec.Cert.TTL = tpo.Spec.TTL
			cro.Spec.VersionBundle.Version = tpo.Spec.VersionBundle.Version

			fmt.Printf("\n")
			fmt.Printf("cro start\n")
			fmt.Printf("%#v\n", cro)
			fmt.Printf("cro end\n")
			fmt.Printf("\n")
		}

		// Create CRO in Kubernetes API.
		{
			_, err := clientSet.CoreV1alpha1().CertConfigs(tpo.Namespace).Get(cro.Name, apismetav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				_, err := clientSet.CoreV1alpha1().CertConfigs(tpo.Namespace).Create(cro)
				if err != nil {
					logger.Log("error", fmt.Sprintf("%#v", err))
					return
				}
			} else if err != nil {
				logger.Log("error", fmt.Sprintf("%#v", err))
				return
			}

			time.Sleep(1 * time.Second)
		}
	}

	logger.Log("debug", "end TPR migration")
}
