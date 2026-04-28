#!/bin/bash
# shellcheck disable=SC2034

set -e

# Colors
GREEN="\e[32m"
DEFAULT="\e[0m"
RED="\e[31m"

SECONDS=0
function display_usage() {
  echo -e "Provisioning Script v0.1\n"
  echo -e "Usage: $0  <argument>"
}

function setup_bash_profile() {
  echo -e "${GREEN}Setting up profile...${DEFAULT}"
  cat <<-EOF >> ~/.bashrc

# Custom PS1
export PS1='\u@\h \W $ '

# User Aliases
alias ls='ls -al'

# Terraform
alias tfi="terraform init"
alias tff="terraform fmt -recursive"
alias tfaa="terraform apply -auto-approve"
alias tfdaa="terraform destroy -auto-approve"

# AWS
alias awsc="aws configure"
alias awswho="aws sts get-caller-identity"

# Extend PATH
export PATH=$PATH:/scriptbin:/.tfenv/bin
EOF
}

function install_tools() {
  # VSCode
  echo -e "${GREEN}Installing VSCode...${DEFAULT}"
  sudo rpm --import https://packages.microsoft.com/keys/microsoft.asc
  cat <<EOF | sudo tee /etc/yum.repos.d/microsoft.repo
[vscode]
name=Visual Studio Code
baseurl=https://packages.microsoft.com/yumrepos/vscode
enabled=1
gpgcheck=1
gpgkey=https://packages.microsoft.com/keys/microsoft.asc
EOF
  sudo dnf install code -y
}

function install_cloud_tools() {
  curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip awscliv2.zip
    sudo ./aws/install
}

generate_gpg_key() {
  gpg --full-generate-key
}

case $1 in
  --vscode)
  install_tools
  ;;
  --set-bash-profile)
  setup_bash_profile
  ;;
  --install-cloud-tools)
  install_cloud_tools
  ;;
  *)
  display_usage
  ;;
esac

echo -e "${GREEN}Elapsed time: $SECONDS seconds.${DEFAULT}"
