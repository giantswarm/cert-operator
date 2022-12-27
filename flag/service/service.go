package service

import (
	"github.com/giantswarm/operatorkit/v8/pkg/flag/service/kubernetes"

	"github.com/giantswarm/cert-operator/v3/flag/service/app"
	"github.com/giantswarm/cert-operator/v3/flag/service/crd"
	"github.com/giantswarm/cert-operator/v3/flag/service/resource"
	"github.com/giantswarm/cert-operator/v3/flag/service/vault"
)

type Service struct {
	App        app.App
	CRD        crd.CRD
	Kubernetes kubernetes.Kubernetes
	Resource   resource.Resource
	Vault      vault.Vault
}
