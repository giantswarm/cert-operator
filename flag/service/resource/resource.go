package resource

import (
	"github.com/giantswarm/cert-operator/v2/flag/service/resource/vaultcrt"
)

type Resource struct {
	VaultCrt vaultcrt.VaultCrt
}
