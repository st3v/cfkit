#!/bin/bash 

set -x

docker build -t st3v/cfkit-test .
docker images
docker ps -a
docker run --rm -e CF_API -e CF_DOMAIN -e CF_USERNAME -e CF_PASSWORD -e CF_ORG -e CF_RABBIT_SERVICE_NAME -e CF_RABBIT_SERVICE_PLAN st3v/cfkit-test /bin/sh -c "ginkgo -r --race -randomizeAllSpecs -cover $@"
docker rmi -f st3v/cfkit-test
docker ps -a
docker images
