#!/bin/bash
eval $(minikube docker-env)
export GOOS=linux 
go build .
mv keysaas ./artifacts/keysaas
docker build -t keysaas-operator:0.5.0 ./artifacts/
rm ./artifacts/keysaas