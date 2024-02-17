#!/bin/bash

# This script is used to install a GitLab Runner on a Linux machine

echo "Installing dependencies"
sudo apt -y update
sudo apt -y install ca-certificates curl gnupg lsb-release
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable | sudo tee /etc/apt/sources.list.d/docker.list >/dev/null
curl -fsSL https://deb.nodesource.com/setup_$NODE_VERSION.x | sudo -E bash -
sudo apt -y update
sudo apt -y install docker-ce docker-ce-cli containerd.io jq git unzip
sudo systemctl enable containerd.service
sudo service docker start

echo "Creating user \"$RUNNER_USER\""
sudo useradd -m $RUNNER_USER
sudo usermod -aG docker $RUNNER_USER

echo "Creating action runner directory"
mkdir -p /opt/github-runner

filename="actions-runner.tar.gz"
echo "Downloading action runner from $RUNNER_DOWNLOAD_URL to $filename"
curl -L $RUNNER_DOWNLOAD_URL -o $filename

echo "Extracting action runner"
tar xzf $filename -C /opt/github-runner
echo "Delete tar file"
rm -rf $filename

echo "Install Script finished"