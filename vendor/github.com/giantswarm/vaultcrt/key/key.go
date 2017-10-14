package key

import (
	"fmt"

	vaultrolekey "github.com/giantswarm/vaultrole/key"
)

func IssuePath(ID string, organizations []string) string {
	return fmt.Sprintf("pki-%s/issue/%s", ID, vaultrolekey.RoleName(ID, organizations))
}
