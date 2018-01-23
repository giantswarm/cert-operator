[![CircleCI](https://circleci.com/gh/giantswarm/cert-operator.svg?style=shield)](https://circleci.com/gh/giantswarm/cert-operator) [![Docker Repository on Quay](https://quay.io/repository/giantswarm/cert-operator/status "Docker Repository on Quay")](https://quay.io/repository/giantswarm/cert-operator)

# cert-operator

Cert Operator creates/configure/manages certificates for Kubernetes clusters
running on Giantnetes.


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
