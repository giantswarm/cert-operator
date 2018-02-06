package vaultcrt

import (
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/vaultcrt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/cert-operator/service/certconfig/v2/key"
)

const (
	// AllowSubDomains defines whether to allow the generated root CA of the PKI
	// backend to allow sub domains as common names.
	AllowSubDomains = true
	Name            = "vaultcrtv2"
	// ConfigHashAnnotation is the annotation key used to track the hash
	// representation of the cert config. This is used to identify changes of the
	// config to trigger renewals.
	ConfigHashAnnotation = "cert.giantswarm.io/config-hash"
	// UpdateTimestampAnnotation is the annotation key used to track the last
	// update timestamp of certificates contained in the Kubernetes secrets.
	UpdateTimestampAnnotation = "giantswarm.io/update-timestamp"
	// UpdateTimestampLayout is the time layout used to format and parse the
	// update timestamps tracked in the annotations of the Kubernetes secrets.
	UpdateTimestampLayout = "2006-01-02T15:04:05.000000Z"
)

type Config struct {
	CurrentTimeFactory func() time.Time
	K8sClient          kubernetes.Interface
	Logger             micrologger.Logger
	VaultCrt           vaultcrt.Interface

	ExpirationThreshold time.Duration
	Namespace           string
}

func DefaultConfig() Config {
	return Config{
		CurrentTimeFactory: nil,
		K8sClient:          nil,
		Logger:             nil,
		VaultCrt:           nil,

		ExpirationThreshold: 0,
		Namespace:           "",
	}
}

type Resource struct {
	currentTimeFactory func() time.Time
	k8sClient          kubernetes.Interface
	logger             micrologger.Logger
	vaultCrt           vaultcrt.Interface

	expirationThreshold time.Duration
	namespace           string
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.CurrentTimeFactory == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.CurrentTimeFactory must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.VaultCrt == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultCrt must not be empty")
	}

	if config.ExpirationThreshold == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.ExpirationThreshold must not be empty")
	}
	if config.Namespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Namespace must not be empty")
	}

	r := &Resource{
		currentTimeFactory: config.CurrentTimeFactory,
		k8sClient:          config.K8sClient,
		logger: config.Logger.With(
			"resource", Name,
		),
		vaultCrt: config.VaultCrt,

		expirationThreshold: config.ExpirationThreshold,
		namespace:           config.Namespace,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) Underlying() framework.Resource {
	return r
}

func (r *Resource) issueCertificate(customObject v1alpha1.CertConfig) (string, string, string, error) {
	c := vaultcrt.CreateConfig{
		AltNames:      key.AltNames(customObject),
		CommonName:    key.CommonName(customObject),
		ID:            key.ClusterID(customObject),
		IPSANs:        key.IPSANs(customObject),
		Organizations: key.Organizations(customObject),
		TTL:           key.CrtTTL(customObject),
	}
	result, err := r.vaultCrt.Create(c)
	if err != nil {
		return "", "", "", microerror.Mask(err)
	}

	return result.CA, result.Crt, result.Key, nil
}

func toSecret(v interface{}) (*apiv1.Secret, error) {
	if v == nil {
		return nil, nil
	}

	secret, ok := v.(*apiv1.Secret)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &apiv1.Secret{}, v)
	}

	return secret, nil
}
