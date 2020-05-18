module github.com/giantswarm/cert-operator

go 1.13

require (
	github.com/giantswarm/apiextensions v0.2.0
	github.com/giantswarm/certs v0.2.0
	github.com/giantswarm/exporterkit v0.2.0
	github.com/giantswarm/k8sclient v0.2.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.2.0
	github.com/giantswarm/microkit v0.2.1
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/operatorkit v0.2.0
	github.com/giantswarm/vaultcrt v0.2.0
	github.com/giantswarm/vaultpki v0.2.0
	github.com/giantswarm/vaultrole v0.2.0
	github.com/giantswarm/versionbundle v0.2.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/prometheus/client_golang v1.6.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/cluster-api v0.3.2
	sigs.k8s.io/controller-runtime v0.5.1
)
