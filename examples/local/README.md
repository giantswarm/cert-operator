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


## Preparing Templates

All yaml files in this directory are templates. Before proceeding this guide
all placeholders must be replaced with sensible values.

- *CLUSTER_NAME* - Cluster name to be created by [aws-operator].
- *COMMON_DOMAIN* - Cluster name to be created by [aws-operator].
- *VAULT_HOST* - When using Vault service from `vault.yaml` `VAULT_HOST` should
  be `vault`. See Vault Setup section below.
- *VAULT_TOKEN* - It must match across the Vault service and the operator
  deployment flags.

Below is handy snippet than can be used to make that painless. It works in bash and zsh.

```bash
for f in *.tmpl.yaml; do
    sed \
        -e 's/${CLUSTER_NAME}/example-cluster/g' \
        -e 's/${COMMON_DOMAIN}/company.com/g' \
        -e 's/${VAULT_HOST}/vault/g' \
        -e 's/${VAULT_TOKEN}/secret_sauce/g' \
        ./$f > ./${f%.tmpl.yaml}.yaml
done
```

- Note: Single quotes are intentional. Strings like `${CLUSTER_NAME}` shouldn't
  be interpolated. These are placeholders in the template files.


## Vault Setup

The operator needs a connection to Vault (currently v0.6.4 is supported) and to
the Kubernetes API. For development running Vault in dev mode is fine.

Steps below are optional. It's OK to use a different Vault instance accessible
from the operator pod. Remember to set `VAULT_HOST` during templates
preparation accordingly.

```bash
kubectl apply -f ./vault.yaml
```

If you are using minikube you can access vault with:
```bash
export VAULT_TOKEN=<VAULT_TOKEN>
export VAULT_ADDR=$(minikube service vault --url)
vault status
```


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
docker build -t quay.io/giantswarm/cert-operator:local-dev .

# Optional. Restart running operator after image update.
# Does nothing when the operator is not deployed.
#kubectl delete pod -l app=cert-operator-local
```


## Operator Startup

```bash
kubectl apply -f ./deployment.yaml
```


## Creating Certificates CustomObjects

```bash
for f in *-cert.yaml; do
    kubectl create -f ./$f
done
```

The certificates are issued using Vault and stored as K8s secrets.

```bash
kubectl get secret -l clusterID=CLUSTER_NAME
```


## Cleaning Up

Delete the certificate custom objects and the deployment.

```bash
kubectl delete certificate -l clusterID=CLUSTER_NAME
kubectl delete -f ./deployment.yaml

# Optinal. Only when Vault was set up.
kubectl delete -f ./vault.yaml
```
