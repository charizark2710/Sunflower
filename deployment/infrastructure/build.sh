#!/bin/bash

# Usage: ./build.sh [dev|prod] [ip]
export COMPOSE_BAKE=true

if [ "$1" == "prod" ]; then
    docker swarm init --advertise-addr "$2"
else
    docker swarm init --advertise-addr 127.0.0.1
fi

docker network create --driver overlay --attachable pi4
docker-compose -f docker-compose.infra.yml --env-file ./../.env up -d --build
