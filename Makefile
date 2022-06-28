up: ## bring app and all its supporting services up
	docker-compose up --build

down: ## shutdown app and all its supporting services
	docker-compose down

test: ## build docker image [standalone] & assign it latest
	go test -v ./...

build-n-tag: ## build docker image [standalone] & assign it latest
	docker build --tag bhardwaz007/geolocation:latest .

push: ## push the built image to hub.docker
	docker push bhardwaz007/geolocation:latest

create-network: ## create a separate network bridge for the app
	docker network create geolocation-net

create-postgres: ## create a postgres databse in the network
	docker run --name postgres -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres --net geolocation-net postgres:alpine

create-postgres-explorer: ## create a postgres-explorer databse in the network
	docker run -p 5050:80 -e "PGADMIN_DEFAULT_EMAIL=someone@postgres.com" -e "PGADMIN_DEFAULT_PASSWORD=postgres" -d --net geolocation-net 

ingest: ## pulls docker image [if not already present] & runs ingest
	docker run -it \
    -p 8080:8080 \
    --net geolocation-net \
    -v ${HOME}/geolocation/data_dump.csv:/root/data_dump.csv \
    -v ${HOME}/geolocation/.env:/root/.env \
    bhardwaz007/geolocation:latest ingest

serve: ## pulls docker image [if not already present] & runs serve
	docker run -it \
    -p 8080:8080 \
    --net geolocation-net \
    -v ${HOME}/geolocation/data_dump.csv:/root/data_dump.csv \
    -v ${HOME}/geolocation/.env:/root/.env \
    bhardwaz007/geolocation:latest serve

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help