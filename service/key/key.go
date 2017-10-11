package key

import (
	"fmt"
	"strings"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
)

func ClusterID(customObject certificatetpr.CustomObject) string {
	return customObject.Spec.ClusterID
}

func ToCustomObject(v interface{}) (certificatetpr.CustomObject, error) {
	customObjectPointer, ok := v.(*certificatetpr.CustomObject)
	if !ok {
		return certificatetpr.CustomObject{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &certificatetpr.CustomObject{}, v)
	}
	customObject := *customObjectPointer

	return customObject, nil
}

func VaultAllowedDomains(customObject certificatetpr.CustomObject, commonNameFormat string) string {
	commonName := VaultCommonName(customObject, commonNameFormat)
	domains := append([]string{commonName}, VaultAltNames(customObject)...)
	return strings.Join(domains, ",")
}

func VaultAltNames(customObject certificatetpr.CustomObject) []string {
	return customObject.Spec.AltNames
}

func VaultAllowBareDomains(customObject certificatetpr.CustomObject) string {
	return customObject.AllowBareDomains
}

func VaultCommonName(customObject certificatetpr.CustomObject, commonNameFormat string) string {
	return fmt.Sprintf(commonNameFormat, ClusterID(customObject))
}

func VaultListMountsPath(customObject certificatetpr.CustomObject) string {
	return fmt.Sprintf("pki-%s", ClusterID(customObject))
}

func VaultListRolesPath(customObject certificatetpr.CustomObject) string {
	return fmt.Sprintf("pki-%s/roles/", ClusterID(customObject))
}

func VaultMountPKIPath(customObject certificatetpr.CustomObject) string {
	return fmt.Sprintf("pki-%s", ClusterID(customObject))
}

func VaultReadCAPath(customObject certificatetpr.CustomObject) string {
	return fmt.Sprintf("pki-%s/cert/ca", ClusterID(customObject))
}

func VaultRoleName(customObject certificatetpr.CustomObject) string {
	return fmt.Sprintf("role-%s", ClusterID(customObject))
}

func VaultWriteCAPath(customObject certificatetpr.CustomObject) string {
	return fmt.Sprintf("pki-%s/root/generate/internal", ClusterID(customObject))
}

func VaultWriteRolePath(customObject certificatetpr.CustomObject) string {
	return fmt.Sprintf("pki-%s/roles/%s", ClusterID(customObject), VaultRoleName(customObject))
}
