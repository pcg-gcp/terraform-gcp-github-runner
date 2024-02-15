#! /bin/bash

nerdctl build . -t europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane && nerdctl push europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane
gcloud run deploy control-plane \
	--set-env-vars="PROJECT_ID=cw-td-sandbox,ZONE=europe-west1-b" \
	--set-env-vars="IMAGE_PATH=projects/debian-cloud/global/images/family/debian-12, MACHINE_TYPE=e2-micro" \
	--max-instances=2 --no-allow-unauthenticated --region=europe-west1 \
	--service-account=control-plane-sa@cw-td-sandbox.iam.gserviceaccount.com \
	--image=europe-docker.pkg.dev/cw-td-sandbox/docker-repo/control-plane
