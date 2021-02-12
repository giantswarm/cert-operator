package key

import (
	"crypto/sha1" // nolint: gosec
	"encoding/json"
	"fmt"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/certs/v3/pkg/certs"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/cert-operator/pkg/label"
)

const (
	CAID  = "ca"
	CrtID = "crt"
	KeyID = "key"
)

func AllowBareDomains(customObject v1alpha1.CertConfig) bool {
	return customObject.Spec.Cert.AllowBareDomains
}

func AltNames(customObject v1alpha1.CertConfig) []string {
	return customObject.Spec.Cert.AltNames
}

func ClusterComponent(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.Cert.ClusterComponent
}

func ClusterID(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.Cert.ClusterID
}

func CommonName(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.Cert.CommonName
}

func ClusterNamespace(customObject v1alpha1.CertConfig) string {
	return ClusterID(customObject)
}

func CrtTTL(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.Cert.TTL
}

// nolint: gosec
func CustomObjectHash(customObject v1alpha1.CertConfig) (string, error) {
	b, err := json.Marshal(customObject.Spec.Cert)
	if err != nil {
		return "", microerror.Mask(err)
	}

	h := sha1.New()
	if _, err := h.Write(b); err != nil {
		return "", err
	}
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs), nil
}

func IPSANs(customObject v1alpha1.CertConfig) []string {
	return customObject.Spec.Cert.IPSANs
}

func IsDeleted(customObject v1alpha1.CertConfig) bool {
	return customObject.GetDeletionTimestamp() != nil
}

func OperatorVersion(customObject v1alpha1.CertConfig) string {
	return customObject.Labels[label.OperatorVersion]
}

func Organizations(customObject v1alpha1.CertConfig) []string {
	a := []string{customObject.Spec.Cert.ClusterComponent}
	return append(a, customObject.Spec.Cert.Organizations...)
}

func RoleTTL(customObject v1alpha1.CertConfig) string {
	return customObject.Spec.Cert.TTL
}

func SecretName(customObject v1alpha1.CertConfig) string {
	cert := certs.Cert(customObject.Spec.Cert.ClusterComponent)
	return certs.K8sName(ClusterID(customObject), cert)
}

func SecretLabels(customObject v1alpha1.CertConfig) map[string]string {
	cert := certs.Cert(customObject.Spec.Cert.ClusterComponent)
	return certs.K8sLabels(ClusterID(customObject), cert)
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
