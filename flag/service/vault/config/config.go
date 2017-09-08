package config

import (
	"github.com/giantswarm/cert-operator/flag/service/vault/config/pki"
)

type Config struct {
	Address string
	Token   string

	PKI pki.PKI
}
