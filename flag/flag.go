package flag

import (
	"github.com/giantswarm/microkit/flag"

	"github.com/giantswarm/cert-operator/flag/kubernetes"
	"github.com/giantswarm/cert-operator/flag/vault"
)

type Flag struct {
	Kubernetes kubernetes.Kubernetes
	Vault      vault.Vault
}

func New() *Flag {
	f := &Flag{}
	flag.Init(f)
	return f
}
