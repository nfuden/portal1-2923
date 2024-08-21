#!/bin/sh


export GLOO_EDGE_PORTAL_VERSION="1.4.2"
export GLOO_EDGE_PORTAL_HELM_VALUES_FILE="./install/ge-portal-helm-values.yaml"
kubectl label namespace default  discovery.solo.io/function_discovery=enabled

kubectl create namespace httpbin --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace portal-env --dry-run=client -o yaml | kubectl apply -f -

printf "\nInstalling Gloo Edge Portal $GLOO_EDGE_PORTAL_VERSION ...\n"
helm upgrade --install gloo-portal gloo-portal/gloo-portal --namespace gloo-portal --create-namespace -f $GLOO_EDGE_PORTAL_HELM_VALUES_FILE --version $GLOO_EDGE_PORTAL_VERSION
printf "\n"
