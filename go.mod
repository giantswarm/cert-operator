module github.com/giantswarm/cert-operator

go 1.13

require (
	github.com/giantswarm/apiextensions v0.4.1
	github.com/giantswarm/certs v0.2.0
	github.com/giantswarm/exporterkit v0.2.0
	github.com/giantswarm/k8sclient v0.2.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/microkit v0.2.2
	github.com/giantswarm/micrologger v0.5.0
	github.com/giantswarm/operatorkit v0.2.1
	github.com/giantswarm/vaultcrt v0.2.0
	github.com/giantswarm/vaultpki v0.2.0
	github.com/giantswarm/vaultrole v0.2.0
	github.com/giantswarm/versionbundle v0.2.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/prometheus/client_golang v1.9.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.17.9
	k8s.io/apimachinery v0.17.9
	k8s.io/client-go v0.17.9
	sigs.k8s.io/cluster-api v0.3.11
	sigs.k8s.io/controller-runtime v0.5.11
)
