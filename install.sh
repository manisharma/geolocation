#!/bin/sh

function line() { 
    echo '\n'
    num="${2:-100}"; 
    printf -- "-%.0s" $(seq 1 $num); 
    echo '\n'
}

line
# create network
echo 'creating network'
docker network create geolocation-network
line
echo 'instantiating postgress'
docker run --name postgres -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres --net geolocation-net postgres:alpine
line
echo 'initialising postgres explorer'
docker run -p 5050:80 -e "PGADMIN_DEFAULT_EMAIL=someone@postgres.com" -e "PGADMIN_DEFAULT_PASSWORD=postgres" -d --net geolocation-net 
line
echo 'ingesting location data'
	docker run -it \
    -p 8080:8080 \
    --net geolocation-net \
    -v ${HOME}/geolocation/data_dump.csv:/root/data_dump.csv \
    -v ${HOME}/geolocation/.env:/root/.env \
    bhardwaz007/geolocation:latest ingest
    serve: ## pulls docker image [if not already present] & runs serve
line
echo 'serving'
    docker run -it \
    -p 8080:8080 \
    --net geolocation-net \
    -v ${HOME}/geolocation/data_dump.csv:/root/data_dump.csv \
    -v ${HOME}/geolocation/.env:/root/.env \
    bhardwaz007/geolocation:latest serve