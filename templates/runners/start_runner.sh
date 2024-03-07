#!/bin/bash

echo "Starting runner..."

echo "Download runner config"
config = $(curl -v "http://metadata.google.internal/computeMetadata/v1/instance/github_runner_config" -H "Metadata-Flavor: Google")
instance_name = $(curl -v "http://metadata.google.internal/computeMetadata/v1/instance/name" -H "Metadata-Flavor: Google")
runner_user = $(curl -v "http://metadata.google.internal/computeMetadata/v1/instance/runner_user" -H "Metadata-Flavor: Google")
runner_dir = $(curl -v "http://metadata.google.internal/computeMetadata/v1/instance/runner_dir" -H "Metadata-Flavor: Google")

echo "github_runner_config: $github_runner_config"

cd $runner_dir

echo "Configuring runner"
./config.sh --work _work --name "ghr-$instance_name" $config

echo "Installing as service"
./svc.sh install $runner_user

echo "Starting service"
./svc.sh start

echo "Runner started. Exiting"
