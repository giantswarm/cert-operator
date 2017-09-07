#!/bin/sh

set -ex

. ./env.sh

for f in *.tmpl.yaml; do
    sed \
        -e 's|${CLUSTER_NAME}|'"${CLUSTER_NAME}"'|g' \
        -e 's|${COMMON_DOMAIN}|'"${COMMON_DOMAIN}"'|g' \
        -e 's|${VAULT_HOST}|'"${VAULT_HOST}"'|g' \
        -e 's|${VAULT_TOKEN}|'"${VAULT_TOKEN}"'|g' \
        ./$f > ./${f%.tmpl.yaml}.yaml
done

kubectl apply -f ./vault.yaml
eval $(minikube docker-env)
(
    cd ../..
    CGO_ENABLED=0 GOOS=linux go build .
    docker build -t quay.io/giantswarm/cert-operator:local-dev .
)

kubectl apply -f ./deployment.yaml

for f in *-cert.yaml; do
    kubectl create -f ./$f
done
