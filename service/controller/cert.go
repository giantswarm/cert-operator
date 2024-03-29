package controller

import (
	"context"
	"fmt"
	"time"

	corev1alpha1 "github.com/giantswarm/apiextensions/v6/pkg/apis/core/v1alpha1"
	providerv1alpha1 "github.com/giantswarm/apiextensions/v6/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/controller"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultpki"
	"github.com/giantswarm/vaultpki/key"
	"github.com/giantswarm/vaultrole"
	vaultapi "github.com/hashicorp/vault/api"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/cert-operator/v3/pkg/label"
)

type CertConfig struct {
	K8sClient   k8sclient.Interface
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client

	UniqueApp           bool
	CATTL               string
	CRDLabelSelector    string
	CommonNameFormat    string
	ExpirationThreshold time.Duration
	Namespace           string
	ProjectName         string
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

	var resources []resource.Interface
	{
		c := ResourceSetConfig{
			CtrlClient:  config.K8sClient.CtrlClient(),
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

		resources, err = NewResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var selector labels.Selector
	{
		if config.UniqueApp {
			selector = label.KubeconfigSelector()
		} else {
			selector = label.AppVersionSelector()
		}

		config.Logger.Debugf(context.Background(), "Watching CertConfigs with selector %v", selector)
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Name:      config.ProjectName,
			Resources: resources,
			Selector:  selector,
			NewRuntimeObjectFunc: func() client.Object {
				return new(corev1alpha1.CertConfig)
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

	err = cleanupPKIBackends(config.Logger, config.K8sClient, vaultPKI)
	if err != nil {
		// We don't want a cleanup error to prevent the controller from starting.
		config.Logger.Log("level", "error", "message", "failed to clean up PKI backends", "stack", fmt.Sprintf("%#v", err))
	}

	return c, nil
}

func cleanupPKIBackends(logger micrologger.Logger, k8sClient k8sclient.Interface, vaultPKI vaultpki.Interface) error {
	mounts, err := vaultPKI.ListBackends()
	if err != nil {
		return microerror.Mask(err)
	}

	logger.Log("level", "debug", "message", "cleaning up PKI backends")

	var latestError *error

	for k := range mounts {
		id := key.ClusterIDFromMountPath(k)

		exists, err := tenantClusterExists(k8sClient, id)
		if err != nil {
			return microerror.Mask(err)
		}

		if !exists {
			logger.Log("level", "debug", "message", fmt.Sprintf("deleting PKI backend for Tenant Cluster %#q", id))

			{
				err := k8sClient.CtrlClient().DeleteAllOf(
					context.Background(),
					&corev1alpha1.CertConfig{},
					client.MatchingLabels{label.Cluster: id},
				)
				if errors.IsNotFound(err) {
					// fall through
				} else if err != nil {
					latestError = &err
					logger.Log("level", "error", "message", fmt.Sprintf("error deleting certconfigs for Tenant Cluster %#q", id))
					continue
				}
			}

			{
				err := vaultPKI.DeleteBackend(id)
				if err != nil {
					latestError = &err
					logger.Log("level", "error", "message", fmt.Sprintf("error deleting PKI backend for Tenant Cluster %#q", id))
					continue
				}
			}

			logger.Log("level", "debug", "message", fmt.Sprintf("deleted PKI backend for Tenant Cluster %#q", id))
		}
	}

	if latestError != nil {
		return microerror.Mask(*latestError)
	}

	logger.Log("level", "debug", "message", "cleaned up PKI backends")

	return nil
}

func tenantClusterExists(k8sClient k8sclient.Interface, id string) (bool, error) {
	var err error

	// We need to check for Node Pools clusters. These adhere to CAPI and do not
	// have any AWSConfig CR anymore.
	{
		crs := &capi.ClusterList{}

		var labelSelector client.MatchingLabels
		{
			labelSelector = make(map[string]string)
			labelSelector[label.Cluster] = id
		}

		err := k8sClient.CtrlClient().List(context.Background(), crs, labelSelector)
		if errors.IsNotFound(err) {
			// fall through
		} else if IsNoKind(err) {
			// fall through
		} else if err != nil {
			return false, microerror.Mask(err)
		} else if len(crs.Items) < 1 {
			// fall through
		} else {
			return true, nil
		}
	}

	// We need to check for the legacy KVMConfig CRs on KVM environments.
	{
		err = k8sClient.CtrlClient().Get(context.Background(), types.NamespacedName{Name: id, Namespace: corev1.NamespaceDefault}, &providerv1alpha1.KVMConfig{})
		if errors.IsNotFound(err) {
			// fall through
		} else if IsNoKind(err) {
			// fall through
		} else if err != nil {
			return false, microerror.Mask(err)
		} else {
			return true, nil
		}
	}

	return false, nil
}
