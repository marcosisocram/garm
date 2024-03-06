CMD_DOCKER_COMPOSE=docker compose
CMD_DOCKER_COMPOSE_UP=$(CMD_DOCKER_COMPOSE) up
LOG_DOCKER_COMPOSE=docker-compose.logs

ENVS=INIT_TABLES=true PORT=8180

PROGRAM_NAME=rinha-go

build:
	go build -o $(PROGRAM_NAME)
grun:
	go run .
brun: clean build
	$(ENVS) ./$(PROGRAM_NAME)
clean:
	rm $(PROGRAM_NAME)

dup:
	$(CMD_DOCKER_COMPOSE_UP)
dupd:
	$(CMD_DOCKER_COMPOSE_UP) --force-recreate -V -d
ddown:
	$(CMD_DOCKER_COMPOSE) down
logs: dupd
	$(CMD_DOCKER_COMPOSE) logs -t --follow > $(LOG_DOCKER_COMPOSE)

