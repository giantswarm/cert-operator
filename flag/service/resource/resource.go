package resource

import (
	"github.com/giantswarm/cert-operator/v3/flag/service/resource/vaultcrt"
)

type Resource struct {
	VaultCrt vaultcrt.VaultCrt
}
