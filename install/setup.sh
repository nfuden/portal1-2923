#!/bin/sh

pushd ../

# Create httpbin namespace if it does not exist yet
kubectl create namespace httpbin --dry-run=client -o yaml | kubectl apply -f -

# Create portal-env namespace if it does not exist yet
kubectl create namespace portal-env --dry-run=client -o yaml | kubectl apply -f -

printf "\nInstalling K8S Services & Deployments ...\n"
kubectl apply -f apis/httpbin/httpbin.yaml
kubectl apply -f apis/petstore/petstore.yaml
printf "\n"

printf "\nInstalling Upstreams ...\n"
kubectl apply -f upstreams/httpbin-httpbin-8000.yaml
printf "\n"

printf "\nInstalling APIDocs ...\n"
kubectl apply -f apis/httpbin/httpbin-schema-apidoc.yaml
kubectl apply -f apis/petstore/petstore-schema-apidoc.yaml
printf "\n"

printf "\nInstalling API Products ...\n"
kubectl apply -f apis/httpbin/httpbin-product-apiproduct.yaml
kubectl apply -f apis/httpbin/httpbin2-product-apiproduct.yaml
kubectl apply -f apis/petstore/petstore-apiproduct.yaml
kubectl apply -f apis/petstore/petstore-get-apiproduct.yaml
kubectl apply -f apis/petstore/petstore-post-apiproduct.yaml
printf "\n"

printf "\nInstalling Portal Environments ...\n"
kubectl apply -f environment/dev-environment.yaml
kubectl apply -f environment/dev-2-environment.yaml
printf "\n"

printf "\nInstalling RouteTables ...\n"
kubectl apply -f routetables/petstore-product.v1.yaml
kubectl apply -f routetables/dev.petstore-product.v1.yaml
kubectl apply -f routetables/dev-2.petstore-product.v1.yaml
printf "\n"

printf "\nInstalling Users and Groups ...\n"
kubectl apply -f users-groups/developers-group-users-secret.yaml
printf "\n"

printf "\nInstalling Portal ...\n"
kubectl apply -f portal/petstore-portal-portal.yaml
printf "\n"
 
popd