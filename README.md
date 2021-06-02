[![CircleCI](https://circleci.com/gh/giantswarm/cert-operator.svg?style=shield)](https://circleci.com/gh/giantswarm/cert-operator) [![Docker Repository on Quay](https://quay.io/repository/giantswarm/cert-operator/status "Docker Repository on Quay")](https://quay.io/repository/giantswarm/cert-operator)

# cert-operator

Cert Operator creates, configures, and manages certificates for Kubernetes clusters
running on the Giant Swarm platform.

Most of the functionality currently provided by this project is now supported natively by Kubernetes' Cluster API (CAPI). As we move more platform functionality to use CAPI workflows, this project will eventually be deprecated.

## About

`cert-operator` is responsible for provisioning certificates used by components of the Giant Swarm platform. It reconciles [`CertConfig` Custom Resources](https://docs.giantswarm.io/ui-api/management-api/crd/certconfigs.core.giantswarm.io/) (CRs) and configures Hashicorp `vault` accordingly. For a given `CertConfig`, `cert-operator` ensures:

- `vault` is accessible
- the necessary `vault` PKI backend has been created
- a root CA for the associated workload cluster has been created using the PKI backend

Secrets are then created in the management cluster containing the certificates, signed by the root CA, used for establishing connections with and within the workload cluster.
Currently, `cert-operator` handles creation of kubeconfigs for workload cluster access for the following components:

- the Giant Swarm API
- app-operator
- aws/azure/kvm-operator
- calico
- etcd
- node-operator
- Prometheus

### Compatibility

| provider   | cert-operator | cluster-operator |
| -----------|:---------:|:---------:|
| AWS        | < 1.0.0*  | < 3.6.1   |
| AWS        | >= 1.0.1  | >= 3.6.1  |
| all others | >= 1.0.1  | >= 0.24.1 |
| all others | < 1.0.1   | < 0.24.0* |

\* cert-operator v1.0.0 and cluster-operator v0.24.0 have known issues. Use v1.0.1 or v0.24.1 instead.

Prior to version 1.0.0, `cert-operator` reconciled based on the `spec.versionBundle.version` field of the `CertConfig` CR.

In version 1.0.0 and later, the CR field is ignored, and the operator reconciles `CertConfig`s which have the `cert-operator.giantswarm.io/version` label set to the operator's version.

In a typical pre-CAPI Giant Swarm release, `cluster-operator` creates the `CertConfig`s necessary for each cluster. `cluster-operator` prior to version 3.6.1 (AWS) and 0.24.0 (Azure and KVM) did not set the appropriate label and still used the older hardcoded `versionBundle`. The two methods are not compatible.

## Prerequisites

## Getting Project

Download the latest release:
https://github.com/giantswarm/cert-operator/releases/latest

Clone the git repository: https://github.com/giantswarm/cert-operator.git

Download the latest docker image from here:
https://quay.io/repository/giantswarm/cert-operator


### How to build


#### Dependencies

- [github.com/giantswarm/microkit](https://github.com/giantswarm/microkit)


#### Building the standard way

```
go build github.com/giantswarm/cert-operator
```


## Running cert-operator

See [this guide][examples-local].

[examples-local]: https://github.com/giantswarm/cert-operator/blob/master/examples/README.md


## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/cert-operator/issues)


## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the
contribution workflow as well as reporting bugs.


## License

cert-operator is under the Apache 2.0 license. See the [LICENSE](LICENSE) file
for details.


## Credit
- https://golang.org
- https://github.com/giantswarm/microkit


### Secrets

The cert-operator is deployed via Kubernetes.

Here the plain Vault token has to be inserted.

```
service:
  vault:
    config:
      token: 'TODO'
```

Here the base64 representation of the data structure above has to be inserted.

```
apiVersion: v1
kind: Secret
metadata:
  name: cert-operator-secret
  namespace: giantswarm
type: Opaque
data:
  secret.yaml: 'TODO'
```

To create the secret manually do this.

```
kubectl create -f ./path/to/secret.yaml
```
