#!/bin/bash

export GOOS=linux; go build .
mv moodle ./artifacts/moodle
docker build -t moodle-operator:0.4.5 ./artifacts/
rm ./artifacts/moodle



