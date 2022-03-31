package service

import (
	"github.com/giantswarm/operatorkit/v7/pkg/flag/service/kubernetes"

	"github.com/giantswarm/cert-operator/flag/service/crd"
	"github.com/giantswarm/cert-operator/flag/service/resource"
	"github.com/giantswarm/cert-operator/flag/service/vault"
)

type Service struct {
	CRD        crd.CRD
	Kubernetes kubernetes.Kubernetes
	Resource   resource.Resource
	Vault      vault.Vault
}
