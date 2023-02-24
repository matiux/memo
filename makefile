# Setup ————————————————————————————————————————————————————————————————————————————————————————————————————————————————
PROJECT_PREFIX=memo

# Static ———————————————————————————————————————————————————————————————————————————————————————————————————————————————
.DEFAULT_GOAL := help
NODEJS_IMAGE=$(PROJECT_PREFIX)-nodejs
DB-IMAGE=$(PROJECT_PREFIX)-db
PROJECT_NAME=$(shell basename $$(pwd) | tr '[:upper:]' '[:lower:]')

# Docker conf ——————————————————————————————————————————————————————————————————————————————————————————————————————————

ifeq ($(wildcard ./docker/docker-compose.override.yml),)
	COMPOSE_OVERRIDE=
else
	COMPOSE_OVERRIDE=-f ./docker/docker-compose.override.yml
endif

COMPOSE=docker compose --file ./docker/docker-compose.yml $(COMPOSE_OVERRIDE) -p $(PROJECT_NAME)
COMPOSE_RUN=$(COMPOSE) run --rm
COMPOSE_EXEC=$(COMPOSE) exec

# Docker commands ——————————————————————————————————————————————————————————————————————————————————————————————————————
.PHONY: up
up: ## Up dei container
	$(COMPOSE) up $$ARG

.PHONY: upd
upd: ## Up dei container in modalità demone
	$(COMPOSE) up -d $$ARG

.PHONY: down
down: ## Down dei container
	$(COMPOSE) down $$ARG

.PHONY: purge
purge: ## Down dei container e pulizia di immagini e volumi
	$(COMPOSE) down --rmi=all --volumes --remove-orphans

.PHONY: log
log: ## Log dei container docker
	$(COMPOSE) logs -f

.PHONY: ps
ps: ## Lista dei container
	@$(COMPOSE) ps

.PHONY: compose
compose: ## Wrapper a docker compose
	@$(COMPOSE) $$ARG

# Commitlint commands ——————————————————————————————————————————————————————————————————————————————————————————————————

.PHONY: conventional
conventional: ## chiama conventional commit per validare l'ultimo commit message
	$(COMPOSE_RUN) -T $(NODEJS_IMAGE) commitlint -e --from=HEAD -V

# Database
.PHONY: prepare-db
prepare-db: ## prepara il database
	$(COMPOSE_EXEC) $(DB-IMAGE) /bin/bash -c "mysql -u root -proot < /docker-entrypoint-initdb.d/sql-schema.sql"

.PHONY: help
help:	## Show this help
	@grep -hE '^[A-Za-z0-9_ \-]*?:.*##.*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'