package pki

import (
	"github.com/giantswarm/cert-operator/v2/flag/service/vault/config/pki/ca"
	"github.com/giantswarm/cert-operator/v2/flag/service/vault/config/pki/commonname"
)

type PKI struct {
	CA         ca.CA
	CommonName commonname.CommonName
}
