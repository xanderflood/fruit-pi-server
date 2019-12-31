.PHONY: build-docker build-local docker local

build-docker:
	CGO_ENABLED=0 GOOS=linux go build -o build/api/api ./cmd/api
	docker build build/api -t xanderflood/fruit-pi-server:local

build-local:
	go build -o build/api/api ./cmd/api

# use this to prep: docker run -p 5432:5432 -e POSTGRES_USER=fruit_pi_server -e POSTGRES_DB=fruit_pi_server postgres
docker: build-docker
	docker run --publish 8000:8000 --env-file .docker.env xanderflood/fruit-pi-server:local

# use this to prep: docker run -p 5432:5432 -e POSTGRES_USER=fruit_pi_server -e POSTGRES_DB=fruit_pi_server postgres
local: build-local
	cd build/api && godotenv -f ../../.env ./api
