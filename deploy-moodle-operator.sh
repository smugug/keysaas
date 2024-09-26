#!/bin/bash

echo "Deploying Keysaas Operator"

MINIKUBE_IP=`minikube ip`

echo "MINIKUBE IP:$MINIKUBE_IP"

rm -f artifacts/deploy-keysaas-operator.yaml

sed "s/MINIKUBE_IP/$MINIKUBE_IP/g" artifacts/deploy-keysaas-operator-minikube.yaml > artifacts/deploy-keysaas-operator.yaml

kubectl create -f artifacts/deploy-keysaas-operator.yaml

echo "Done."

echo "You can now create Keysaas instances as follows:"
echo "kubectl apply -f artifacts/keysaas1.yaml"
echo "kubectl describe keysaases keysaas1"
