module github.com/giantswarm/cert-operator

go 1.13

require (
	github.com/giantswarm/apiextensions/v3 v3.18.2
	github.com/giantswarm/certs/v3 v3.1.1
	github.com/giantswarm/exporterkit v0.2.1
	github.com/giantswarm/k8sclient/v5 v5.0.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/microkit v0.2.2
	github.com/giantswarm/micrologger v0.5.0
	github.com/giantswarm/operatorkit/v4 v4.2.0
	github.com/giantswarm/vaultcrt v0.2.0
	github.com/giantswarm/vaultpki v0.2.0
	github.com/giantswarm/vaultrole v0.2.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/prometheus/client_golang v1.9.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.18.15
	k8s.io/apimachinery v0.18.15
	k8s.io/client-go v0.18.15
	sigs.k8s.io/cluster-api v0.3.13
	sigs.k8s.io/controller-runtime v0.6.4
)

replace sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.13-gs
