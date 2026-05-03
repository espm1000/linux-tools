#!/bin/bash
# shellcheck source=/dev/null
# shellcheck disable=SC2320 # Ignore exit commands from an 'echo'
# shellcheck disable=SC2154 # Var referenced but not assigned
# shellcheck disable=SC2086 # Double quotes

set -e

# Colors
GREEN="\e[32m"
DEFAULT="\e[0m"
RED="\e[31m"

DEP_LIST="dos2unix sudo ca-certificates curl gnupg2 git cifs-utils net-tools"
DEV_TOOLS="build-essential checkinstall libz-dev dh-make libssl-dev devscripts"

# Prompt for what to run
# selection is a var in the read line.
echo "Options:"
echo "1 --- Install Basic Deps"
echo "2 --- Install Dev Tools"
echo "3 --- Install Docker-CE"
echo "4 --- Install kubectl"
echo "5 --- Install TFenv"
echo "6 --- Install GCloud SDK"
read -rp "Make selection: " selection

install_basic_deps() {
echo "PS1='${debian_chroot:+($debian_chroot)}\u@\h:/\W\$ '" >> ~/.bashrc
if [ "$USER" != 'root' ]; then
  echo -e "${RED}You need to run this as root.${DEFAULT}"
  exit 1
else
  apt update && apt upgrade -y
  for dep in ${DEP_LIST}; do
	  apt install $dep -y
  done
  echo "alias ls='ls -al'" >> ~/.bashrc
  echo "alias ls='ls -al'" >> /home/nick/.bashrc
  if [ -e /home/"$USER"/.vimrc ]; then
    echo "vimrc file found, assuming already configured.  Skipping."
  else
    echo "set nocompatible" >> ~/.vimrc
    echo "set backspace=indent,eol,start" >> ~/.vimrc
    echo "set nocompatible" >> /home/nick/.vimrc
    echo "set backspace=indent,eol,start" >> /home/nick/.vimrc
  fi
  export PATH=$PATH:/usr/sbin
  usermod -aG sudo nick
  cp initProvisioning.sh /home/nick
  echo -e "${RED}Please reboot.${DEFAULT}"
fi
}

install_dev_tools() {
  sudo apt update
  echo -e "${GREEN}Installing Dev Tools${DEFAULT}"
  sudo apt install ${DEV_TOOLS} -y
}

install_docker() {
  # check if docker is installed
    if ! command -v docker &> /dev/null
    then
      echo -e "${RED}Docker not found, fixing.${DEFAULT}"
    if [[ ! -e '/etc/apt/keyrings/docker.gpg' ]]; then
      echo -e "${RED}No keyring found, fixing.${DEFAULT}"
      sudo install -m 0755 -d /etc/apt/keyrings
      curl -fsSL https://download.docker.com/linux/debian/gpg \
      | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
      sudo chmod a+r /etc/apt/keyrings/docker.gpg
      # Add repo
      echo \
      "deb [arch=""$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
      "$(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
      sudo apt update
    else
      echo -e "${GREEN}Keyring found, skipping.${DEFAULT}"
    fi
    # Install docker binaries
    sudo apt install -y \
    docker-ce \
    docker-ce-cli \
    containerd.io \
    docker-buildx-plugin \
    docker-compose-plugin
  else
    echo -e "${GREEN}Docker found. Exiting.${DEFAULT}"
    exit 0
  fi
}

install_kubectl() {
  FINAL_PATH="/usr/local/bin/"
  echo -e "${GREEN}Installing KubeCTL...${DEFAULT}"
  if ! command -v curl &> /dev/null; then
    echo -e "${RED}cURL not installed...please install first. ${DEFAULT}"
    exit 1
  fi
  curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" || echo "download failed"
  curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl.sha256" || echo "downloading checksum failed"
  echo -e "${GREEN}Checking sum...${DEFAULT}"
  echo "$(cat kubectl.sha256) kubectl" | sha256sum --check || echo "failed checksum"
  chmod +x ./kubectl && sudo mv ./kubectl ${FINAL_PATH} || echo -e "${RED}failed to copy file${DEFAULT}"
  rm -f kubectl*
  echo -e "${GREEN}Installation complete.  Update PATH to include '${FINAL_PATH}'"
}

install_tfenv() {
  if ! command -v git &> /dev/null; then
    echo -e "${RED}Git not installed...please install first. ${DEFAULT}"
    exit 1
  fi
  git clone --depth=1 https://github.com/tfutils/tfenv.git ~/.tfenv || echo "install failed"
  echo -e "${GREEN}Install complete. Set the environment to include '~/.tfenv/bin' ${DEFAULT}"
}

install_google_cli() {
  google_key_url="https://packages.cloud.google.com/apt/doc/apt-key.gpg"
  sudo apt update && sudo apt install -y \
    apt-transport-https \
    ca-certificates \
    gnupg \
    curl

    echo -e "${GREEN}Adding Google GPG Key${DEFAULT}"
    curl $google_key_url | sudo gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
    echo -e "${GREEN}Done.${DEFAULT}"
    
    echo -e "${GREEN}Installing GCloud SDK${DEFAULT}"
    sudo apt update && sudo apt install -y \
      google-cloud-cli

}

if [ ${selection} = 1 ]; then
  install_basic_deps
elif [ ${selection} = 2 ]; then
  install_dev_tools
elif [ ${selection} = 3 ]; then
  install_docker
elif [ ${selection} = 4 ]; then
  install_kubectl
elif [ ${selection} = 5 ]; then
  install_tfenv
elif [ ${selection} = 6 ]; then
  install_google_cli
else
  exit 1
fi
