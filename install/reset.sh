#!/bin/sh

#
# Reset the environment.
#


#!/bin/sh

pushd ../

printf "\nDelete Portal ...\n"
kubectl delete -f portal/petstore-portal-portal.yaml
printf "\n"

printf "\nDelete Users and Groups ...\n"
kubectl delete -f users-groups/developers-group-users-secret.yaml
printf "\n"

printf "\nDelete RouteTables ...\n"
kubectl delete -f routetables/petstore-product.v1.yaml
kubectl delete -f routetables/dev.petstore-product.v1.yaml
kubectl delete -f routetables/dev-2.petstore-product.v1.yaml
printf "\n"

printf "\nDelete Portal Environments ...\n"
kubectl delete -f environment/dev-environment.yaml
kubectl delete -f environment/dev-2-environment.yaml
printf "\n"

printf "\nDelete API Products ...\n"
kubectl delete -f apis/httpbin/httpbin-product-apiproduct.yaml
kubectl delete -f apis/httpbin/httpbin2-product-apiproduct.yaml
kubectl delete -f apis/petstore/petstore-apiproduct.yaml
kubectl delete -f apis/petstore/petstore-get-apiproduct.yaml
kubectl delete -f apis/petstore/petstore-post-apiproduct.yaml
printf "\n"

printf "\nDelete APIDocs ...\n"
kubectl delete -f apis/httpbin/httpbin-schema-apidoc.yaml
kubectl delete -f apis/petstore/petstore-schema-apidoc.yaml
printf "\n"

printf "\nDelete Upstreams ...\n"
kubectl delete -f upstreams/httpbin-httpbin-8000.yaml
printf "\n"

printf "\nInstalling K8S Services & Deployments ...\n"
kubectl delete -f apis/httpbin/httpbin.yaml
kubectl delete -f apis/petstore/petstore.yaml
printf "\n"

# Create httpbin namespace if it does not exist yet
kubectl delete namespace httpbin

# Create portal-env namespace if it does not exist yet
kubectl delete namespace portal-env

 
popd