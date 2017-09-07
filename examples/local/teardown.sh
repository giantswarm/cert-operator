#!/bin/sh

set -ex

. ./env.sh

kubectl delete certificate -l clusterID=${CLUSTER_NAME}
kubectl delete -f ./deployment.yaml
kubectl delete -f ./vault.yaml
