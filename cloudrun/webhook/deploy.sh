#! /bin/bash

nerdctl build . -t europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook && nerdctl push europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook
gcloud run deploy webhook \
	--max-instances=2 --allow-unauthenticated --region=europe-west1 \
	--set-env-vars="PROJECT_ID=cw-td-sandbox,TASK_QUEUE_PATH=projects/cw-td-sandbox/locations/europe-west1/queues/github-events" \
	--set-env-vars="CONTROL_PLANE_URL=https://control-plane-qljv23sbqa-ew.a.run.app,DEBUG=true" \
	--set-env-vars="INVOKER_SERVICE_ACCOUNT=control-plane-invoker@cw-td-sandbox.iam.gserviceaccount.com" \
	--set-secrets="WEBHOOK_SECRET_KEY=github_webhook_secret:latest" \
	--service-account=webhook-sa@cw-td-sandbox.iam.gserviceaccount.com \
	--image=europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook
