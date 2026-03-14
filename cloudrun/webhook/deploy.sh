#! /bin/bash

docker build . -t europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook:latest && docker push europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook:latest
gcloud run deploy ghr-webhook \
  --region=europe-west1 \
  --project=cw-td-sandbox \
  --image=europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook:latest
