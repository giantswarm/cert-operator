package config

import (
	"github.com/giantswarm/cert-operator/flag/service/vault/config/certificate"
	"github.com/giantswarm/cert-operator/flag/service/vault/config/pki"
)

type Config struct {
	Address string
	Token   string

	Certificate certificate.Certificate
	PKI         pki.PKI
}
