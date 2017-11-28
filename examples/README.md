# Running cert-operator Locally

**Note:** This should only be used for testing and development. See the
[/kubernetes/][kubernetes-dir] directory and [Secrets][secrets-doc] for
a production ready configuration.

[kubernetes-dir]: https://github.com/giantswarm/cert-operator/tree/master/kubernetes
[secrests-doc]: https://github.com/giantswarm/cert-operator#secrets

This guide explains how to get running cert-operator locally. For example on
minikube. Certificates created here are meant to be used by [aws-operator].

All commands are assumed to be run from `examples/local` directory.

[aws-operator]: https://github.com/giantswarm/aws-operator

## Cluster-Local Docker Image

The operator needs a connection to the K8s API. The simplest approach is to run
as a deployment and use the "in cluster" configuration.

In that case the Docker image needs to be accessible from the K8s cluster
running the operator. For Minikube `eval $(minikube docker-env)` before `docker
build`, see [reusing the Docker daemon] for details.

[reusing the docker daemon]: https://github.com/kubernetes/minikube/blob/master/docs/reusing_the_docker_daemon.md

```bash
# Optional. Only when using Minikube.
eval $(minikube docker-env)

# From the root of the project, where the Dockerfile resides
CGO_ENABLED=0 GOOS=linux go build github.com/giantswarm/cert-operator
docker build -t quay.io/giantswarm/cert-operator:local-lab .

# Optional. Restart running operator after image update.
# Does nothing when the operator is not deployed.
#kubectl delete pod -l app=cert-operator-local
```

## Deploying the lab charts

The lab consist of three Helm charts, `cert-operator-lab-chart`, which sets up cert-operator,
`cert-resource-lab-chart`, which puts in place the required certificates and `vaultlab-chart`,
which installs Vault in dev mode. For installing the latter two you need the [Helm registry plugin](https://github.com/app-registry/appr-helm-plugin)

With a working Helm installation they can be created from the project's root with:

```bash
$ helm registry install quay.io/giantswarm/vaultlab-chart:stable -- \
                        -n vault \
                        --set vaultToken=myToken

$ helm install -n cert-operator-lab \
               --set imageTag=local-lab \
               --set vaultToken=myToken \
               --set commonDomain=mydomain.io \
               ./examples/cert-operator-lab-chart/ --wait

helm registry install quay.io/giantswarm/cert-resource-lab-chart:stable -- \
                      -n cert-resource-lab \
                      --set commonDomain=mydomain.io \
                      --set clusterName=test-cluster
```

The certificates are issued using Vault and stored as K8s secrets.

```bash
kubectl get secret -l clusterID=test-cluster # or the actual value of `clusterName`
```

`cert-operator-lab-chart` accepts the following configuration parameters:
* `commonDomain` - Domain to be used by [aws-operator].
* `vaultHost` - Defaults to `vault` for the local setup.
* `vaultToken` - It must match across the Vault service and the operator deployment flags.
* `imageTag` - Tag of the cert-operator image to be used, by default `local-dev` to use a locally created
image.

`cert-resource-lab-chart` is also configurable with `clusterName` and `commonDomain` (the latter should match the value
used in `cert-operator-lab-chart`).


You can specify different values of the configuration parameters changing the `values.yaml` file on each
chart directory or specifying them on the install command:
```bash
$ helm install -n cert-operator-lab --set clusterName=my-cluste-name ./cert-operator-lab-chart/ --wait
```

## Cleaning Up

Delete the cert-operator and certificates lab releases:

```bash
$ helm delete cert-resource-lab --purge
$ helm delete cert-operator-lab --purge
```
