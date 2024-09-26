#!/bin/bash

echo "Deploying MySQL"

kubectl create -f artifacts/keysaas-mysql.yaml

