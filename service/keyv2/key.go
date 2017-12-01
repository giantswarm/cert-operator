package keyv2

import (
	"fmt"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
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
