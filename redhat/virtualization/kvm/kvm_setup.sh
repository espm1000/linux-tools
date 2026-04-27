#!/bin/bash
# shellcheck disable=all

set -e

# Colors
GREEN="\e[32m"
DEFAULT="\e[0m"
RED="\e[31m"

if [[ $1 == '' ]]; then
    echo "Usage: $0 <vm_name> <vm_ram> <vm_cpu>"
    exit 1
fi

install_virt() {
    # Install DNF Group

    sudo dnf install @virtualization

    # Start service
    sudo systemctl start libvirtd
    sudo systemctl status libvirtd

    # Virtual KVM Modules are loaded

    lsmod | grep kvm
}

# Set up Virtual Machine

# Windows

create_windows11_desktop() {
    # Default location is /var/lib/libvirt/images
    VM_NAME=$1
    VM_RAM=$2
    VM_CPU=$3
    VM_DISK_SIZE=$4

    VM_PATH=/home/$USER/virtual_machines/kvm_images/${VM_NAME}.qcow2
    if [[ ! -e ${VM_PATH} ]]; then
    echo -e "${RED}Image file not found, creating...${DEFAULT}"
    echo sudo qemu-img create -f qcow2 ${VM_PATH} ${VM_DISK_SIZE}
    sudo qemu-img create -f qcow2 ${VM_PATH} ${VM_DISK_SIZE}
    echo -e "${GREEN}Disk created successfully.${DEFAULT}"
    else
    echo -e "${GREEN}Disk file found.${DEFAULT}"
    fi
    
    # Create Storage using QCOW2

    VM_CDROM_PATH='/home/nick/Downloads/ISO/windows11/Win11_24H2_English_x64.iso'

    echo sudo virt-install --name ${VM_NAME} \
        --description "'KVM Workstation'" \
        --ram ${VM_RAM} \
        --vcpus ${VM_CPU} \
        --disk ${VM_PATH},size=${VM_DISK_SIZE} \
        --cdrom ${VM_CDROM_PATH}

    sudo virt-install --name ${VM_NAME} \
        --description "'KVM Workstation'" \
        --ram ${VM_RAM} \
        --vcpus ${VM_CPU} \
        --disk ${VM_PATH},size=64 \
        --cdrom ${VM_CDROM_PATH} \
        --osinfo win11 
        
}

main() {
    create_windows11_desktop "$@"
}

main "$@"
