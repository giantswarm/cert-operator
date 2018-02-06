package key

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
)

const (
	CAID  = "ca"
	CrtID = "crt"
	KeyID = "key"
)

const (
	// ComponentLabel is the label used in the secret to identify a cluster
	// component.
	ComponentLabel string = "clusterComponent"
	// ClusterIDLabel is the label used in the secret to identify a cluster.
	ClusterIDLabel string = "clusterID"
)

func AllowBareDomains(customObject v1alpha1.CertConfig) bool {
	return customObject.Spec.Cert.AllowBareDomains
}

func AltNames(customObject v1alpha1.CertConfig) []string {
	return customObject.Spec.Cert.AltNames
}

func ClusterID(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.Cert.ClusterID
}

func ClusterComponent(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.Cert.ClusterComponent
}

func CommonName(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.Cert.CommonName
}

func CrtTTL(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.Cert.TTL
}

func CustomObjectHash(customObject v1alpha1.CertConfig) (string, error) {
	b, err := json.Marshal(customObject.Spec.Cert)
	if err != nil {
		return "", microerror.Mask(err)
	}

	h := sha1.New()
	h.Write(b)
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs), nil
}

func IPSANs(customObject v1alpha1.CertConfig) []string {
	return customObject.Spec.Cert.IPSANs
}

func Organizations(customObject v1alpha1.CertConfig) []string {
	a := []string{customObject.Spec.Cert.ClusterComponent}
	return append(a, customObject.Spec.Cert.Organizations...)
}

func SecretName(customObject v1alpha1.CertConfig) string {
	return fmt.Sprintf("%s-%s", ClusterID(customObject), ClusterComponent(customObject))
}

func RoleTTL(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.Cert.TTL
}

func ToCustomObject(v interface{}) (v1alpha1.CertConfig, error) {
	customObjectPointer, ok := v.(*v1alpha1.CertConfig)
	if !ok {
		return v1alpha1.CertConfig{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha1.CertConfig{}, v)
	}
	customObject := *customObjectPointer

	return customObject, nil
}

func VersionBundleVersion(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.VersionBundle.Version
}
