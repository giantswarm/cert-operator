package pki

import (
	"github.com/giantswarm/cert-operator/flag/service/vault/config/pki/ca"
	"github.com/giantswarm/cert-operator/flag/service/vault/config/pki/commonname"
)

type PKI struct {
	CA         ca.CA
	CommonName commonname.CommonName
}
