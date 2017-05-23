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
https://hub.docker.com/r/giantswarm/cert-operator/


### How to build


#### Dependencies

- [github.com/giantswarm/microkit](https://github.com/giantswarm/microkit)
- [github.com/giantswarm/certificatetpr](https://github.com/giantswarm/certificatetpr)


#### Building the standard way

```
go build github.com/giantswarm/cert-operator
```


## Running cert-operator

The operator needs a connection to Vault (currently v0.6.4 is supported) and to
the Kubernetes API. For development running Vault in dev mode is fine.


### Setup

- The operator needs to connect to a Vault server. See
  [examples/vault.yaml](https://github.com/giantswarm/cert-operator/blob/master/examples/vault.yaml)
  for running Vault as a deployment with a ClusterIP service.
- The cert-operator binary needs to be built into a docker image and tagged as
  `quay.io/giantswarm/cert-operator:local-dev`. The current pod need to be
  deleted for changes to apply.

```
GOOS=linux go build github.com/giantswarm/cert-operator \
  && docker build -t quay.io/giantswarm/cert-operator:local-dev . \
  && kubectl delete pod -l app=cert-operator-local
```

- The docker image needs to be accessible from the k8s cluster. For Minikube
  see [reusing the docker daemon](https://github.com/kubernetes/minikube/blob/master/docs/reusing_the_docker_daemon.md).
- The operator also needs a connection to the K8s API. The simplest approach is
  to run as a deployment and use the "in cluster" configuration.

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: cert-operator-local
  namespace: default
  labels:
    app: cert-operator-local
spec:
  replicas: 1
  strategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: cert-operator-local
    spec:
      volumes:
      containers:
      - name: cert-operator
        image: quay.io/giantswarm/cert-operator:local-dev
        imagePullPolicy: IfNotPresent
        ports:
        - name: http
          containerPort: 8000
        args:
        - daemon
        - --service.vault.config.address=http://YOUR_VAULT_HOST:8200
        - --service.vault.config.token=YOUR_TOKEN
        - --service.vault.config.pki.ca.ttl=1440h
        - --service.vault.config.pki.commonname.format=%s.g8s.aws.giantswarm.io
```

- Note: Edit YOUR_VAULT_HOST to point at your Vault endpoint.
- Note: This should only be used for development. See the
  [/kubernetes/](https://github.com/giantswarm/cert-operator/tree/master/kubernetes)
  directory and [Secrets](https://github.com/giantswarm/cert-operator#secrets)
  for a production ready configuration.


### Creating TPOs (Third Party Objects)

- The [/examples/](https://github.com/giantswarm/cert-operator/tree/master/examples) directory contains a set of certificatetpr resources designed
to work with the [example cluster](https://github.com/giantswarm/aws-operator/blob/master/examples/cluster.yml) in the `aws-operator`.

```
for i in examples/*-cert.yaml; do kubectl create -f $i; done
```

- The certificates are issued using Vault and stored as k8s secrets.

```
kubectl get secret -l clusterID=example-cluster
```


### Cleaning up

- Delete the certificate TPOs and the deployment.

```
kubectl delete certificate -l clusterID=example-cluster
kubectl delete deployment cert-operator-local
```


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
  secret.yml: 'TODO'
```

To create the secret manually do this.

```
kubectl create -f ./path/to/secret.yml
```
