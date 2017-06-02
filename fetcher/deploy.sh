#!/usr/bin/env bash

# Convert to lowercase and replace space with
POD_NAME=$(echo $1 | awk '{print tolower($0)}' | sed -e "s/ /-/g")

cat k8s.jobs.yaml | sed -e "s/{{SPECIES_NAME}}/$1/g" | sed -e "s/{{POD_NAME}}/$POD_NAME/g" | kubectl create -f -