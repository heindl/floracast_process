#!/usr/bin/env bash

cat k8s.jobs.yaml | sed -e "s/{{SPECIES_NAME}}/$1/g" | kubectl create -f -