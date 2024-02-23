#! /bin/bash

nerdctl build . -t europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook:latest && nerdctl push europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook:latest
gcloud run deploy ghr-webhook \
	--region=europe-west1 \
	--image=europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook:latest
