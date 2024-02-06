#! /bin/bash

nerdctl build . -t europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook && nerdctl push europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook
gcloud run deploy webhook \
	--max-instances=2 --allow-unauthenticated --region=europe-west1 \
	--set-env-vars="PUBSUB_PROJECT_ID=cw-td-sandbox,DEBUG=true" \
	--set-secrets="WEBHOOK_SECRET_KEY=github_webhook_secret:latest" \
	--service-account=webhook-sa@cw-td-sandbox.iam.gserviceaccount.com \
	--image=europe-docker.pkg.dev/cw-td-sandbox/docker-repo/webhook
