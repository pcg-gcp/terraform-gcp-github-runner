#! /bin/bash

nerdctl build . -t europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane:latest && nerdctl push europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane:latest
gcloud run deploy ghr-control-plane \
  --region=europe-west1 \
  --project=cw-td-sandbox \
  --image=europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane:latest
