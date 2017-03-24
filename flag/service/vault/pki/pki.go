package pki

import (
	"github.com/giantswarm/cert-operator/flag/service/vault/pki/ca"
	"github.com/giantswarm/cert-operator/flag/service/vault/pki/commonname"
)

type PKI struct {
	CA         ca.CA
	CATTL      string
	CommonName commonname.CommonName
}
