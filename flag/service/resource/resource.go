package resource

import (
	"github.com/giantswarm/cert-operator/flag/service/resource/vaultcrt"
)

type Resource struct {
	VaultCrt vaultcrt.VaultCrt
}
