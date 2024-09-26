#!/bin/bash

export GOOS=linux; go build .
mv keysaas ./artifacts/keysaas
docker build -t keysaas-operator:0.4.5 ./artifacts/
rm ./artifacts/keysaas



