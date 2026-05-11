#!/bin/bash
set -euo pipefail

# Usage: ./deploy-swarm.sh [dev|prod] [ip] [--skip-infra]
MODE="${1:-dev}"
IP="${2:-127.0.0.1}"
SKIP_INFRA=false

# Check for the skip flag in any position
for arg in "$@"; do
    if [ "$arg" == "--skip-infra" ]; then
        SKIP_INFRA=true
    fi
done

# 1. Handle Credentials (unchanged)
GenerateRandomPw () {
    if [ -f "./.env" ]; then
        echo "Found existing .env file. Skipping password generation to maintain DB compatibility."
        return
    fi

    echo "Generating fresh production credentials..."
    declare -a keys=("postgres" "broker" "keycloak" "mongo" "api")

    for key in "${keys[@]}"; do
        # Generates a secure 24-character hex string
        local pw=$(openssl rand -hex 12)
        # Correct way to export dynamic variable names
        export "${key}_pw"="$pw"
    done

    if [ -f "./local.deploy.env" ]; then
        envsubst < ./local.deploy.env > ./.env
        echo ".env file created successfully."
    else
        echo "ERROR: local.deploy.env template not found!"
        exit 1
    fi
}

GenerateRandomPw

# 2. Conditional Infra Deployment
if [ "$SKIP_INFRA" = true ]; then
    echo "Checking if infrastructure is already running..."
    if docker stack ls | grep -q "infra"; then
        echo "Infra stack detected. Skipping build/deploy for infrastructure."
    else
        echo "WARNING: --skip-infra was passed but no 'infra' stack found. Deploying anyway..."
        bash ./infrastructure/build.sh "./infrastructure" "$MODE" "$IP"
    fi
else
    echo "Deploying infrastructure..."
    bash ./infrastructure/build.sh "./infrastructure" "$MODE" "$IP"

    sleep 1m && echo "Deploying infrastructure Done"
fi

# Export existing env vars for envsubst
set -a
[ -f "./.env" ] && . ./.env
set +a

echo "Preparing swarm stack configuration..."
envsubst < ./docker-compose-swarm.yml > ./deploy-docker-compose-swarm.yml

# 4. Build and Deploy
echo "Building application services..."
docker-compose -f docker-compose-swarm.yml --env-file ./.env build

echo "Deploying stack: sunflower..."
docker stack deploy --compose-file deploy-docker-compose-swarm.yml sunflower

# 5. Cleanup
rm ./deploy-docker-compose-swarm.yml
echo "Deployment of 'sunflower' complete."

# Scale services
# docker service scale sunflower_db=1
# docker service scale sunflower_keycloak=1
# docker service scale sunflower_rabbitmq=1
# docker service scale sunflower_memcached=1
# docker service scale sunflower_api=1
# docker service scale sunflower_server=1
# docker service scale sunflower_pgadmin4=1
