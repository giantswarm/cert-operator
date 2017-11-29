package keyv1

import (
	"fmt"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
)

func AllowBareDomains(customObject certificatetpr.CustomObject) bool {
	return customObject.Spec.AllowBareDomains
}

func AltNames(customObject certificatetpr.CustomObject) []string {
	return customObject.Spec.AltNames
}

func ClusterID(customObject certificatetpr.CustomObject) string {
	return customObject.Spec.ClusterID
}

func ClusterComponent(customObject certificatetpr.CustomObject) string {
	return customObject.Spec.ClusterComponent
}

func CommonName(customObject certificatetpr.CustomObject) string {
	return customObject.Spec.CommonName
}

func CrtTTL(customObject certificatetpr.CustomObject) string {
	return customObject.Spec.TTL
}

func IPSANs(customObject certificatetpr.CustomObject) []string {
	return customObject.Spec.IPSANs
}

func Organizations(customObject certificatetpr.CustomObject) []string {
	a := []string{customObject.Spec.ClusterComponent}
	return append(a, customObject.Spec.Organizations...)
}

func SecretName(customObject certificatetpr.CustomObject) string {
	return fmt.Sprintf("%s-%s", ClusterID(customObject), ClusterComponent(customObject))
}

func RoleTTL(customObject certificatetpr.CustomObject) string {
	return customObject.Spec.TTL
}

func ToCustomObject(v interface{}) (certificatetpr.CustomObject, error) {
	customObjectPointer, ok := v.(*certificatetpr.CustomObject)
	if !ok {
		return certificatetpr.CustomObject{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &certificatetpr.CustomObject{}, v)
	}
	customObject := *customObjectPointer

	return customObject, nil
}

func VersionBundleVersion(customObject certificatetpr.CustomObject) string {
	return customObject.Spec.VersionBundle.Version
}
