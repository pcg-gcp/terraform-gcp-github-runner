#!/bin/bash

# This script is used to install a GitLab Runner on a Linux machine

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
