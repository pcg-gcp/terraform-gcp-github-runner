#!/bin/bash

# This script is used to install a GitHub Runner on a Linux machine

echo "Installing dependencies"
echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections
sudo apt-get -y update
sudo apt-get -y install ca-certificates curl gnupg lsb-release
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable | sudo tee /etc/apt/sources.list.d/docker.list >/dev/null
curl -fsSL https://deb.nodesource.com/setup_$NODE_VERSION.x | sudo -E bash -
sudo apt-get -y update
sudo apt-get -y install docker-ce docker-ce-cli containerd.io jq git unzip
sudo systemctl enable containerd.service
sudo service docker start

echo "Creating user \"$RUNNER_USER\""
sudo useradd -m $RUNNER_USER
sudo usermod -aG docker $RUNNER_USER

echo "Creating action runner directory"
mkdir -p $RUNNER_DIR

filename="actions-runner.tar.gz"
echo "Downloading action runner from $RUNNER_DOWNLOAD_URL to $filename"
curl -sS -L $RUNNER_DOWNLOAD_URL -o $filename

echo "Extracting action runner"
tar xzf $filename -C $RUNNER_DIR
echo "Delete tar file"
rm -rf $filename

echo "Setting permissions"
chown -R $RUNNER_USER $RUNNER_DIR

echo "Install Script finished"

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
