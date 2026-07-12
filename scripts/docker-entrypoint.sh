#!/usr/bin/env sh
# Docker entrypoint for OpenMuara.
# Initializes a default config on first start, then runs the requested muara
# command.

set -e

CONFIG_DIR="/app/.muara"
CONFIG_PATH="${CONFIG_DIR}/config.yml"

# In containers the server should bind to all interfaces so port mappings work.
export MUARA_SERVER_HOST="${MUARA_SERVER_HOST:-0.0.0.0}"

maybe_init() {
  if [ "$1" = "start" ] && [ ! -f "$CONFIG_PATH" ]; then
    echo "No config found at ${CONFIG_PATH}; initializing defaults..."
    # Container deployments default to hardened mode with a random admin token
    # so that 0.0.0.0 binding passes config validation.
    export MUARA_ADMIN_ENABLED="${MUARA_ADMIN_ENABLED:-true}"
    export MUARA_ADMIN_USERNAME="${MUARA_ADMIN_USERNAME:-admin}"
    if [ -z "${MUARA_ADMIN_TOKEN:-}" ] && [ -z "${MUARA_ADMIN_PASSWORD_HASH:-}" ]; then
      TOKEN="$(head -c 32 /dev/urandom | tr -dc 'a-zA-Z0-9' 2>/dev/null || printf '%s' "$(date +%s%N)")"
      export MUARA_ADMIN_TOKEN="${TOKEN}"
      echo "Generated admin token: ${TOKEN}"
    fi
    export MUARA_HARDENED="${MUARA_HARDENED:-true}"
    muara --config "$CONFIG_PATH" init --defaults
  fi
}

maybe_init "$1"

exec muara --config "$CONFIG_PATH" "$@"
