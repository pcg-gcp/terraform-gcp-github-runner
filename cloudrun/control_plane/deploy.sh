#! /bin/bash

docker build . -t europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane:latest && docker push europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane:latest
gcloud run deploy ghr-control-plane \
  --region=europe-west1 \
  --project=cw-td-sandbox \
  --image=europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane:latest
