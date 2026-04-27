#!/bin/bash
# shellcheck disable=SC2248

# Set desired FQDN

COMPOSE_FILE="${COMPOSE_FILE:-docker-compose-openldap-standalone.yaml}"
PERSIST_LDAP_DATA="${PERSIST_LDAP_DATA:-""}"

echo "Using compose file ${COMPOSE_FILE}"

set_hostname() {
    read -rp "Hostname: " ldap_hostname
    read -rp "Domain: " ldap_domain
    ldap_fqdn="${ldap_hostname}.${ldap_domain}"

    echo "Setting fqdn ${ldap_fqdn}"
    echo "127.0.0.1 ${ldap_fqdn} ${ldap_hostname}" >> /etc/hosts
}

set_firewall_rules() {
    command -v firewall-cmd >> /dev/null
    if [[ $? -eq 1 ]]; then
      echo "firewall-cmd not found"
    else
      echo "firewall-cmd found, (re)adding rules."
      sudo firewall-cmd --add-service=ldap --permanent
      sudo firewall-cmd --add-service=ldaps --permanent
      sudo firewall-cmd --reload
    fi
}

configure_docker_compose() {
    if [[ -z ${PERSIST_LDAP_DATA} ]]; then
      echo "Removing data dir.."
      sudo rm -rf data/
    else
      echo "Persisting datastore..."
    fi

    if [[ ! -f ${COMPOSE_FILE} ]]; then
      echo "docker-compose file missing, fix it."
      exit 1
    else
      read -rp "Admin password: " admin_password
      read -rp "Org: " org
      read -rp "Domain: " domain
      read -rp "Suffix: " suffix
      if [[ -f "${COMPOSE_FILE}.running" ]]; then
        mv "${COMPOSE_FILE}.running" "${COMPOSE_FILE}.backup"
      fi
      cp "${COMPOSE_FILE}" "${COMPOSE_FILE}.running"
      sed -i "s/{adminpassword}/${admin_password}/gI" "${COMPOSE_FILE}.running"
      sed -i "s/{org}/${org}/gI" "${COMPOSE_FILE}.running"
      sed -i "s/{domain}/${domain}/gI" "${COMPOSE_FILE}.running"
      sed -i "s/{suffix}/${suffix}/gI" "${COMPOSE_FILE}.running"
      docker compose -f "${COMPOSE_FILE}.running" up -d
    fi
}

main() {
    set_firewall_rules
    configure_docker_compose
}

main
