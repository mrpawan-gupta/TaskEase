DC=docker-compose
COMPOSE_FILE=docker-compose.yml

logs:
	@if [ -z "$$SERVICE" ]; then \
		$(DC) -f $(COMPOSE_FILE) logs; \
	else \
		$(DC) -f $(COMPOSE_FILE) logs $$SERVICE; \
	fi


build:
	$(DC) -f $(COMPOSE_FILE) build

build-up:
	$(DC) -f $(COMPOSE_FILE) up

start:
	$(DC) -f $(COMPOSE_FILE) up -d server

clean-containers:
	docker container prune -f || true
	docker rm -f $$(docker ps -aq) || true

clean-images:
	docker image prune -af || true
	docker rmi -f $$(docker images -q) || true

clean: clean-containers clean-images