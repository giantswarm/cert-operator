package service

import (
	"github.com/giantswarm/cert-operator/flag/service/kubernetes"
	"github.com/giantswarm/cert-operator/flag/service/resource"
	"github.com/giantswarm/cert-operator/flag/service/vault"
)

type Service struct {
	Kubernetes kubernetes.Kubernetes
	Resource   resource.Resource
	Vault      vault.Vault
}
