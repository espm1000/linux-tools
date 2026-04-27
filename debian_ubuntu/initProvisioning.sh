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
echo "3 --- Install SSH Keys"
echo "4 --- Install VSCode Server"
echo "5 --- Install Virtual Box Tools"
echo "6 --- Install Docker-CE"
echo "7 --- Install Corp SSH"
echo "8 --- Install cryptsetup dependencies"
echo "9 --- Install Kubernetes Repo"
read -rp "Make selection: " selection

PUB_KEY="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDEZ1CVZlhCF4bUhtMRaJob8D1T4UbOQeTroU998Aj1QzeSwEQPkCeHklFpmBDbGOrTbbeh9EaopTp9Fy9lJ7ZEG/a6z57fw35t7HOdhp4jvj0p7gTjkljpjpaM1cSHh4noZy0wjCgge8GIWZvr6b8eil+cZQxDleO91QCHg78HiG+8QL5kcOURfpM2dTecmbbAQWdN4vruRnio9/nGmMrydacYg3qSqmi983rpRmZMuFNwUAqfgfiJ0IlwcR3F29I0YZGvzdYyp4XvhQUwxHvC4Nq4g0tZHW8jkzhs14kJ1m5+TfbdHFHyT88Y3BlMJwsdkRnQUQ6kBYtTZayMGgye/He8I8Lw93EEgzc+/duEBOHlWSJSU918+gUTc0rPchbBlXrs8GJ8KytSE84jLc+SLnfTTSL2osZvgspgR+Ob9yUU4pIQn0s1pR5zco5Ft/5FKyj1PMvvup7l4I2QLlOqxqa9J28mjbshKovUslooIBqBOpFzOPpZSM0XgazMZGaJtNg9oq8/mAL6bfv+3OV5QPP6VEXGVHitphzlWZkYcGtiVj7qtio23rWC7TyvanjeLpIH/8L3540dnJ5p9MbFAavVaQX2HSrGkvyqoHOwl63Zx9U8uCCorRh6oo96e/6wjtG0dWwaJajpW54r2DSgaiKf3lPrDLcNnJO69MLMNQ== naspenwall@dragos.com"
VSCODE_PUB_KEY="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC0WdPhmDAZp8RThdp0Z8K8weF4NYlmQAi8TVAeW5lWdiMWaAJ8C1/dGU0Fzgbz0uDs8Y8Zxtm2ne/VLIJn/P+Cn1ZhpJ4/1oDdKUC6/DbxL3xa8qjMkO3NFCBMxkgAMtJ8tchwZ62MQAhQ7Q7m59NOcSnC6HEd/AY5f4QG28pYwWZRLjCJBbWasCfv3GSpFh4npvDUCONPpRYXyl1E80gpWPESvyB1UrPnpJl0CM9B+RTr/utVOYuPaMcLGXvVE/eKSxAZdvYWnBupu4dqTBPik+FY9svUQcK/v4UJjOd9CTW155YhwiBWisq8BCOxjw/MsGY6WjnG9ONofnmYiitsVBvjfGgjiuZf7gr/MzOGu+QaG9Shm3Lyeb3B5Gb2Wfy1nPOmfhRs7Q1hd0rWjLnHTuM5xrL/4d2RQMYRa/5Uc1/Fpa+YWDIio71UccEpFjDLXib50+44AFdepzvJHt3nFWqkEW4R1HNBcWrZT0N78dcXwpkj/ufjXZekS0Jo9hZZnhS1m1omshJ0+XMv9SDxr7Qelh6zGYLOdS6c2bljLeFs+ptH+Q5A31b7bzaOtPQGi4s6+O5uSSEj1jr1veB+ixmzH29IvhhY3EWyVeRXnkdSb6ROIzGk48ZmTZu9431aBFia1kJP8z3qTXXpmedrrEh9mnFUmHMs01+lQVFXww== nick@LWD179224"

install_basic_deps() {
echo "PS1='${debian_chroot:+($debian_chroot)}\u@\h:/\W\$ '" >> ~/.bashrc
if [ "$USER" != 'root' ]; then
  echo -e "${RED}You need to run this as root.${DEFAULT}"
  exit 1
else
  apt update && apt upgrade -y
  for dep in "${DEP_LIST}"; do
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

install_corp_ssh_key() {
  # Check if id_rsa exists
if [ ! -e "/home/$USER/usrcrt.pfx" ]; then
  echo -e "${RED}PFX file not found, please upload.${DEFAULT}"
else
  echo -e "${GREEN}Installing public corp key for bitbucket.${DEFAULT}"
  openssl pkcs12 -in /home/"$USER"/usrcrt.pfx -nocerts -nodes \
  | openssl rsa | tee /home/"$USER"/.ssh/id_rsa > /dev/null
  echo -e "${GREEN}Done.${DEFAULT}"
fi
}

install_ssh_key() {
  # Lets set the PS1 value now:
echo "PS1='${debian_chroot:+($debian_chroot)}\u@\h:/\W\$ '" >> ~/.bashrc
if [ ! -e /home/"$USER"/.ssh/authorized_keys ]; then
  echo -e "${GREEN}Adding ssh and key files.${DEFAULT}"
  rm -rf .ssh/
  mkdir .ssh && touch /home/"$USER"/.ssh/authorized_keys
  echo "${PUB_KEY}" >> /home/"$USER"/.ssh/authorized_keys
  echo "${VSCODE_PUB_KEY}" >> /home/$USER/.ssh/authorized_keys
else
  echo "${PUB_KEY}" >> /home/"$USER"/.ssh/authorized_keys
  echo "${VSCODE_PUB_KEY}" >> /home/$USER/.ssh/authorized_keys
  if [ ! $? ]; then
    echo -e "${RED}Something went wrong with the .ssh dir or the authorized_keys file.${DEFAULT}"
  fi
fi
}

install_vscode_server() {
  # Install Common Software
  sudo apt install software-properties-common -y
  # Get Repo
  curl -sSL https://packages.microsoft.com/keys/microsoft.asc | sudo apt-key add -
  # Add repo to apt
  sudo add-apt-repository "deb [arch=amd64] https://packages.microsoft.com/repos/vscode stable main"
  # Install VSCode
  sudo apt update && sudo apt install code -y
  if [ $? = 1 ]; then
    echo -e "${RED}Installation failed.  See apt/journal for more info.${DEFAULT}"
  else
    printf "VSCode server install successful.\nPlease run 'code tunnel --accept-server-license-terms' to continue\n"
  fi
}

install_virtualbox_tools() {
  # Check if cdrom is loaded
  file -s /dev/sr0 > /dev/null
    if [ $? = 1 ]; then
      echo "No media mount.  Please mount VirtualBox tools."
      exit 1
    else
      echo "Media found in drive, attempting to install VB Tools."
      sudo mount /dev/cdrom /mnt && cd /mnt
      sudo ./VBoxLinuxAdditions.run
      #sudo mount -t vboxsf FIPS /mnt
    fi
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

install_luks_deps() {
  # git_branch='v2.3.x'
  if [ -e cryptsetup-deps.txt ]; then
    dep_list="cryptsetup-deps.txt"
  else
    echo "Dep list not found.  Please upload."
    exit 1
  fi
  while IFS= read -r pkg;
    do sudo apt install $pkg -y
  done < "$dep_list"
  # Clone repo from Gitlab.com
  if [ -d "/home/$USER/cryptsetup" ]; then
    echo "/home/$USER/cryptsetup dir already exists, assuming cloned.  Skipping."
  else
    if [ -e cryptsetup-v2.3.x.tar.gz ]; then
      tar zxvf cryptsetup-v2.3.x.tar.gz
      cd cryptsetup-v2.3.x/ && ./autogen.sh && \
      ./configure --prefix=/opt/cryptsetup-install
    else
      echo "Tarball not found. Please upload v2.3.x."
    # git clone https://gitlab.com/cryptsetup/cryptsetup.git /home/$USER/cryptsetup
    # Checkout 2.3 branch
    # echo "Change directory to /home/$USER/cryptsetup"
    # echo "Run command: git check $git_branch"
    fi
  fi
}

install_k8s_repo() {
  # Install k8s deps
  k_deps="apt-transport-https ca-certificates curl"
  for pkg in ${k_deps}; do
    sudo apt install $pkg -y
  done
  # Download public key and add to keyring
  curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.29/deb/Release.key \
  | sudo gpg --dearmor -o /etc/apt/keyrings/k8s.gpg
  # Add repo
  echo 'deb [signed-by=/etc/apt/keyrings/k8s.gpg] https://pkgs.k8s.io/core:/stable:/v1.29/deb /' \
  | sudo tee /etc/apt/sources.list.d/k8s.list
  sudo apt-get update
  echo "Done adding K8s repo."

}

if [ ${selection} = 1 ]; then
  install_basic_deps
elif [ ${selection} = 2 ]; then
  install_dev_tools
elif [ ${selection} = 3 ]; then
  install_ssh_key
elif [ ${selection} = 4 ]; then
  install_vscode_server
elif [ ${selection} = 5 ]; then
  install_virtualbox_tools
elif [ ${selection} = 6 ]; then
  install_docker
elif [ ${selection} = 7 ]; then
  install_corp_ssh_key
elif [ ${selection} = 8 ]; then
  install_luks_deps
elif [ ${selection} = 9 ]; then
  install_k8s_repo
else
  exit 1
fi
