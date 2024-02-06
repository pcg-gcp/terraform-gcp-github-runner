#! /bin/bash

nerdctl build . -t europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane && nerdctl push europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook
gcloud run deploy control-plane \
	--max-instances=2 --no-allow-unauthenticated --region=europe-west1 \
	--service-account=control-plane-sa@cw-td-sandbox.iam.gserviceaccount.com \
	--image=europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane
