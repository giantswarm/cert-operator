# Running cert-operator Locally

```
* TODO: ignore examples/local/vault.yaml
* TODO: ignore examples/local/*-cert.yaml
* TODO: link in root README.md
```

This should be used only for testing end development.

This guide explains how to get running cert-operator locally. For example on
minikube. Certificates created here are meant to be used by [aws-operator].

All commands are assumed to be run from `examples/local` directory.

[aws-operator]: https://github.com/giantswarm/aws-operator

## Preparing Templates

```bash
for f in *.tmpl.yaml; do
    name="${f%.tmpl.yaml}.yaml"
    sed -e 's/${CLUSTER_NAME}/example-cluster/g' ./$f > ./$name
    sed -e 's/${COMMON_DOMAIN}/company.com/g' ./$f > ./$name
    sed -e 's/${VAULT_HOST}/vault/g' ./$f > ./$name
    sed -e 's/${VAULT_TOKEN}/secret_sauce/g' ./$f > ./$name
done
```

- Note: Single quotes are intentional. Strings like `${CLUSTER_NAME}` shouldn't
  be interpolated. They are placeholders in template files.
- Note: When using Vault service from `vault.yaml` `VAULT_HOST` should be
  `vault`. See Vault Setup section below.
- Note: `VAULT_TOKEN` value can be arbitrary. It must match across Vault
  service and the operator deployment flags.

## Vault Setup

The operator needs a connection to Vault (currently v0.6.4 is supported) and to
the Kubernetes API. For development running Vault in dev mode is fine.

Steps below are optional. It's OK to use different Vault instance accessible
from the operator pod. Remember to set `VAULT_HOST` during templates
preparation accordingly. 

```bash
kubectl apply -f ./vault.yaml
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

GOOS=linux go build github.com/giantswarm/cert-operator
docker build -t quay.io/giantswarm/cert-operator:local-dev .

# Optional. Restart running operator after image update.
# Does nothing when the operator is not deployed.
#kubectl delete pod -l app=cert-operator-local
```

## Operator Startup

```bash
kubectl apply -f ./deployment.yaml
```

## Creating Certificates ThirdPartyObjects

```bash
for f in *.-cert.yaml; do
    kubectl create -f ./$f
done
```

The certificates are issued using Vault and stored as K8s secrets.

```bash
kubectl get secret -l clusterID=CLUSTER_NAME
```
