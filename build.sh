#!/bin/sh

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o recipe .
docker build -t recepies_api .
docker tag recepies_api:latest proto:6000/recepies_api
