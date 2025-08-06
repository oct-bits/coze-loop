.PHONY: debug fe server sync_db dump_db middleware web down clean python help

# 定义脚本路径
SCRIPTS_DIR := ./scripts
BUILD_FE_SCRIPT := $(SCRIPTS_DIR)/build_fe.sh
BUILD_SERVER_SCRIPT := $(SCRIPTS_DIR)/setup/server.sh
SYNC_DB_SCRIPT := $(SCRIPTS_DIR)/setup/db_migrate_apply.sh
DUMP_DB_SCRIPT := $(SCRIPTS_DIR)/setup/db_migrate_dump.sh
SETUP_DOCKER_SCRIPT := $(SCRIPTS_DIR)/setup/docker.sh
SETUP_PYTHON_SCRIPT := $(SCRIPTS_DIR)/setup/python.sh
COMPOSE_FILE := docker/docker-compose.yml
MYSQL_SCHEMA := ./docker/volumes/mysql/schema.sql
MYSQL_INIT_SQL := ./docker/volumes/mysql/sql_init.sql
ENV_FILE := ./docker/.env
STATIC_DIR := ./bin/resources/static
ES_INDEX_SCHEMA := ./docker/volumes/elasticsearch/es_index_schema
ES_SETUP_SCRIPT := ./docker/volumes/elasticsearch/setup_es.sh

up:
	docker compose -f ./release/deployment/docker/docker-compose.yml --env-file ./release/deployment/docker/.env --profile "*" up

up-redis:
	docker compose -f ./release/deployment/docker/docker-compose.yml --env-file ./release/deployment/docker/.env --profile "redis" up

up-mysql:
	docker compose -f ./release/deployment/docker/docker-compose.yml --env-file ./release/deployment/docker/.env --profile "mysql" up

up-clickhouse:
	docker compose -f ./release/deployment/docker/docker-compose.yml --env-file ./release/deployment/docker/.env --profile "clickhouse" up

up-minio:
	docker compose -f ./release/deployment/docker/docker-compose.yml --env-file ./release/deployment/docker/.env --profile "minio" up

up-rmq:
	docker compose -f ./release/deployment/docker/docker-compose.yml --env-file ./release/deployment/docker/.env --profile "rmq" up

up-nginx:
	docker compose -f ./release/deployment/docker/docker-compose.yml --env-file ./release/deployment/docker/.env --profile "nginx" up

down:
	docker compose -f ./release/deployment/docker/docker-compose.yml ---profile '*' down

down-v:
	docker compose -f ./release/deployment/docker/docker-compose.yml --profile '*' down -v

image:
	@echo "Building coze-loop-app image..."
	GOPROXY=https://go-mod-proxy.byted.org,https://goproxy.cn,https://proxy.golang.org,direct \
	docker build \
		--progress=plain \
		-f ./release/image/Dockerfile \
		-t open-coze-loop-app:latest \
		.
	docker run --rm open-coze-loop-app:latest du -sh /coze-loop/bin
	docker run --rm open-coze-loop-app:latest du -sh /coze-loop/resources
	docker run --rm open-coze-loop-app:latest du -sh /coze-loop/conf
	docker run --rm open-coze-loop-app:latest du -sh /coze-loop

clean-image:
	docker rmi -f coze-loop-app:latest
	docker builder prune --force

into-image:
	docker run -it --rm open-coze-loop-app:latest /bin/bash

clean-all:
	@echo "Stopping containers..."
	@docker ps -aq | xargs -r docker stop

	@echo "Removing containers..."
	@docker ps -aq | xargs -r docker rm -f

	@echo "Removing images..."
	@docker images -aq | xargs -r docker rmi -f

	@echo "Removing volumes..."
	@docker volume ls -q | xargs -r docker volume rm

	@echo "Removing custom networks..."
	@docker network ls | awk '/bridge|host|none/ {next} NR>1 {print $$1}' | xargs -r docker network rm

	@echo "Pruning builder and system..."
	@docker builder prune -a -f
	@docker system prune -a --volumes -f

build_server:
	@echo "Building image..."
	@bash $(BUILD_SERVER_SCRIPT)

sync_db:
	@echo "Syncing database..."
	@docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) --profile mysql-setup up -d

dump_db: dump_sql_schema
	@echo "Dumping database..."
	@bash $(DUMP_DB_SCRIPT)

sql_init:
	@echo "Init sql data..."
	@docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) --profile mysql-setup up -d

middleware:
	@echo "Start middleware docker environment for opencoze app"
	@docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) --profile middleware up -d --wait


web:
	@echo "Start web server in docker"
	@docker compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) --profile '*' up -d --wait

#down:
#	@echo "Stop all docker containers"
#	@docker compose -f $(COMPOSE_FILE) --profile '*' down

clean: down
	@echo "Remove docker containers and volumes data"
	@rm -rf ./docker/data

python:
	@echo "Setting up Python..."
	@bash $(SETUP_PYTHON_SCRIPT)

dump_sql_schema:
	@echo "Dumping mysql schema to $(MYSQL_SCHEMA)..."
	@. $(ENV_FILE); \
	{ echo "SET NAMES utf8mb4;\nCREATE DATABASE IF NOT EXISTS opencoze COLLATE utf8mb4_unicode_ci;"; atlas schema inspect -u $$ATLAS_URL --format "{{ sql . }}" --exclude "atlas_schema_revisions,table_*" | sed 's/CREATE TABLE/CREATE TABLE IF NOT EXISTS/g'; } > $(MYSQL_SCHEMA)
	@sed -I '' -E 's/(\))[[:space:]]+CHARSET utf8mb4/\1 ENGINE=InnoDB CHARSET utf8mb4/' $(MYSQL_SCHEMA)
	@echo "Dumping mysql schema to helm/charts/opencoze/files/mysql ..."
	@cp $(MYSQL_SCHEMA) ./helm/charts/opencoze/files/mysql/
	@cp $(MYSQL_INIT_SQL) ./helm/charts/opencoze/files/mysql/

atlas-hash:
	@echo "Rehash atlas migration files..."
	@(cd ./docker/atlas && atlas migrate hash)

setup_es_index:
	@echo "Setting up Elasticsearch index..."
	@bash $(ES_SETUP_SCRIPT)  --index-dir $(ES_INDEX_SCHEMA) --docker-host false

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  debug            - Start the debug environment."
	@echo "  env              - Setup env file."
	@echo "  fe               - Build the frontend."
	@echo "  server           - Build and run the server binary."
	@echo "  build_server     - Build the server binary."
	@echo "  sync_db          - Sync opencoze_latest_schema.hcl to the database."
	@echo "  dump_db          - Dump the database to opencoze_latest_schema.hcl and migrations files."
	@echo "  sql_init         - Init sql data..."
	@echo "  dump_sql_schema  - Dump the database schema to sql file."
	@echo "  middleware       - Setup middlewares docker environment, but exclude the server app."
	@echo "  web              - Setup web docker environment, include middlewares docker."
	@echo "  down             - Stop the docker containers."
	@echo "  clean            - Stop the docker containers and clean volumes."
	@echo "  python           - Setup python environment."
	@echo "  atlas-hash       - Rehash atlas migration files."
	@echo "  setup_es_index   - Setup elasticsearch index."
	@echo "  help             - Show this help message."
