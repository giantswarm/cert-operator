package key

import (
	"fmt"
)

func CommonName(ID string, commonNameFormat string) string {
	return fmt.Sprintf(commonNameFormat, ID)
}

func ListMountsPath(ID string) string {
	return fmt.Sprintf("pki-%s", ID)
}

func MountPKIPath(ID string) string {
	return fmt.Sprintf("pki-%s", ID)
}

func ReadCAPath(ID string) string {
	return fmt.Sprintf("pki-%s/cert/ca", ID)
}

func WriteCAPath(ID string) string {
	return fmt.Sprintf("pki-%s/root/generate/internal", ID)
}
