package vault

import (
	"github.com/giantswarm/cert-operator/flag/vault/pki"
)

type Vault struct {
	Address string
	Token   string
	PKI     pki.PKI
}
