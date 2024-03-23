#!/bin/bash

echo "Starting runner..."

echo "Download runner config"
instance_name=$(curl -v "http://metadata.google.internal/computeMetadata/v1/instance/name" -H "Metadata-Flavor: Google")
runner_config=$(curl -v "http://metadata.google.internal/computeMetadata/v1/instance/attributes/github_runner_config" -H "Metadata-Flavor: Google")
runner_user=$(curl -v "http://metadata.google.internal/computeMetadata/v1/instance/attributes/runner_user" -H "Metadata-Flavor: Google")
runner_dir=$(curl -v "http://metadata.google.internal/computeMetadata/v1/instance/attributes/runner_dir" -H "Metadata-Flavor: Google")
use_jit_config=$(curl -v "http://metadata.google.internal/computeMetadata/v1/instance/attributes/use_jit_config" -H "Metadata-Flavor: Google")

cd $runner_dir

if [ "$use_jit_config" = "true" ]; then
	echo "JIT config enabled. Starting runner with JIT config"
	runuser -u $runner_user -- ./run.sh --jitconfig $runner_config
else
	echo "Configuring runner"
	runuser -u $runner_user -- ./config.sh --unattended --work _work --name $instance_name $runner_config

	echo "Installing as service"
	./svc.sh install $runner_user

	echo "Starting service"
	./svc.sh start
fi

echo "Runner started. Exiting"
