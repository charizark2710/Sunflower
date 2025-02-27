#!/bin/sh 

GenerateRandomPw () {
    declare -a keys=("postgres" "broker" "keycloak" "mongo")

    for key in "${keys[@]}"; do
        local pw=$(openssl rand -base64 12 | sha1sum | awk '{print $1}')
        export declare ${key}_pw=$pw
    done
    envsubst < ./local.deploy.env > ./.env
}

#Generate random password for production
GenerateRandomPw

# export all variable in .env
set -a
. ./.env

mkdir -p /tmp/deploy/

# substitube all env in docker-compose-swarm.yml and output to deploy-docker-compose-swarm.yml
envsubst < ./docker-compose-swarm.yml > ./deploy-docker-compose-swarm.yml

# deploy to stack sunflower based on docker-compose-swarm.yml
docker stack deploy --compose-file deploy-docker-compose-swarm.yml sunflower

# Scale services
# docker service scale sunflower_db=1
# docker service scale sunflower_keycloak=1
# docker service scale sunflower_rabbitmq=1
# docker service scale sunflower_memcached=1
# docker service scale sunflower_api=1
# docker service scale sunflower_server=1
# docker service scale sunflower_pgadmin4=1

# remove deploy-docker-compose-swarm.yml after being done
rm ./deploy-docker-compose-swarm.yml
