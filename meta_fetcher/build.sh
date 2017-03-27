#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
chmod 700 main
docker build . -t us.gcr.io/phenograph-154419/species-meta-fetcher
rm ./main
gcloud docker -- push us.gcr.io/phenograph-154419/species-meta-fetcher