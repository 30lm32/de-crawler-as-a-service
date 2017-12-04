CONT_CRAWLER_REDIS := crawler-redis
CONT_CRAWLER_GO := crawler-go
TEST_COMPOSE_FILE := docker-compose.test.yml

build-test:
	docker-compose -f $(TEST_COMPOSE_FILE) build

up-test:
	docker-compose -f $(TEST_COMPOSE_FILE) up

stop-test:
	docker-compose -f $(TEST_COMPOSE_FILE) down
	docker-compose -f $(TEST_COMPOSE_FILE) stop

run-test: build-test up-test

build-service:
	docker-compose build

stop-service:
	docker-compose down
	docker-compose stop

up-service:
	docker-compose up

run-service: build-service up-service

restart:
	docker-compose restart

bash_redis:
	docker exec -it $(CONT_CRAWLER_REDIS) /bin/bash

bash_go:
	docker exec -it $(CONT_CRAWLER_GO) /bin/bash
