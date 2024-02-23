#! /bin/bash

nerdctl build . -t europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane && nerdctl push europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane
gcloud run deploy ghr-control-plane \
	--region=europe-west1 \
	--image=europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane
