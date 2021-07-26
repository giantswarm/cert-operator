module github.com/giantswarm/cert-operator

go 1.13

require (
	github.com/giantswarm/apiextensions/v3 v3.27.3
	github.com/giantswarm/certs/v3 v3.1.1
	github.com/giantswarm/exporterkit v0.2.1
	github.com/giantswarm/k8sclient/v5 v5.11.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/microkit v0.2.2
	github.com/giantswarm/micrologger v0.5.0
	github.com/giantswarm/operatorkit/v4 v4.3.1
	github.com/giantswarm/vaultcrt v0.2.0
	github.com/giantswarm/vaultpki v0.2.0
	github.com/giantswarm/vaultrole v0.2.0
	github.com/hashicorp/vault/api v1.1.0
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.18.19
	k8s.io/apimachinery v0.18.19
	k8s.io/client-go v0.18.19
	sigs.k8s.io/cluster-api v0.3.19
	sigs.k8s.io/controller-runtime v0.6.4
)

replace (
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.24+incompatible
	github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.24+incompatible
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.13-gs
)
