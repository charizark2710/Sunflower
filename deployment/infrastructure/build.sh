#!/bin/bash

# Exit immediately if a command exits with a non-zero status
# Treat unset variables as an error
set -euo pipefail

export COMPOSE_BAKE=true
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT_PATH="$BASE_DIR"

# 1. Improved Path Resolution
if [ -n "${1:-}" ]; then
    SCRIPT_PATH="$(cd "$1" 2>/dev/null && pwd || echo "$1")"
fi

# 2. Idempotent Swarm Init
if ! docker info | grep -q "Swarm: active"; then
    echo "Initializing Swarm..."
    ADVERTISE_ADDR="${3:-127.0.0.1}"
    if [ "${2:-}" == "prod" ]; then
        docker swarm init --advertise-addr "$ADVERTISE_ADDR"
    else
        docker swarm init --advertise-addr 127.0.0.1
    fi
else
    echo "Swarm already active, skipping init."
fi

# 3. Idempotent Network Creation
if ! docker network ls | grep -q "pi4"; then
    echo "Creating overlay network pi4..."
    docker network create --driver overlay --attachable pi4
fi

# 4. Build Images
# Note: Using 'docker buildx bake' is the modern way if COMPOSE_BAKE=true is intended
docker-compose -f "$SCRIPT_PATH/docker-compose.infra.yml" --env-file "$BASE_DIR/../.env" build

# 5. Safe Environment Substitution
set -a
[ -f "$BASE_DIR/../.env" ] && . "$BASE_DIR/../.env"
set +a

# Use a temporary file and ensure it's deleted even if the script fails
TMP_DEPLOY_FILE=$(mktemp "$SCRIPT_PATH/deploy-XXXXXX.yml")
trap 'rm -f "$TMP_DEPLOY_FILE"' EXIT

envsubst < "$SCRIPT_PATH/docker-compose.infra.yml" > "$TMP_DEPLOY_FILE"

# 6. Deploy
docker stack deploy -c "$TMP_DEPLOY_FILE" infra