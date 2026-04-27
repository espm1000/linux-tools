#!/bin/bash

set -e
SECONDS=0

# Get Repo Key
echo "Getting repo key..."
sudo rpm --import https://downloads.1password.com/linux/keys/1password.asc

# Create repo file
echo "Creating repo file..."
echo -e "[1password]\nname=1Password Repository\nbaseurl=https://downloads.1password.com/linux/rpm/stable/\$basearch\nenabled=1\ngpgcheck=1\nrepo_gpgcheck=1\ngpgkey=\"https://downloads.1password.com/linux/keys/1password.asc\"" | sudo tee /etc/yum.repos.d/1password.repo

# Install 1Password
echo "Installing 1Password..."
sudo dnf install 1password -y

echo "Done."
echo "Elapsed time: $SECONDS seconds."
