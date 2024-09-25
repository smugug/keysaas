#!/bin/bash
eval $(minikube docker-env)
export GOOS=linux 
go build .
mv moodle ./artifacts/moodle
docker build -t moodle-operator:0.5.0 ./artifacts/
rm ./artifacts/moodle