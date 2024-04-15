.PHONY: all build up down test logs shell clean stop restart rebuild

all: build

build:
	docker-compose -f docker/docker-compose.yml --env-file .env build

up:
	docker-compose -f docker/docker-compose.yml --env-file .env up -d

down:
	docker-compose -f docker/docker-compose.yml --env-file .env down

test:
	docker-compose -f docker/docker-compose.yml --env-file .env run --rm app go test ./... -v

logs:
	docker-compose -f docker/docker-compose.yml --env-file .env logs

shell:
	docker-compose -f docker/docker-compose.yml --env-file .env exec app sh

clean:
	docker system prune -a
	docker volume prune

stop:
	docker-compose -f docker/docker-compose.yml --env-file .env stop

restart:
	docker-compose -f docker/docker-compose.yml --env-file .env restart

rebuild: down build up
