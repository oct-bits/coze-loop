.PHONY: debug fe server sync_db dump_db middleware web down clean python help

KB_NAMESPACE := coze-loop
KB_RELEASE_NAME := coze-loop
KB_DEPLOY_NAME := coze-loop
KB_CHART_PATH :=./release/deployment/helm-chart/charts
KB_UMBRELLA_PATH :=./release/deployment/helm-chart/umbrella

DOCKER_COMPOSE_DIR := ./release/deployment/docker-compose
#demo-mini:
#	kubectl -n ingress-nginx patch svc ingress-nginx-controller -p '{"spec":{"type":"LoadBalancer"}}' service/ingress-nginx-controller patched \
#    kubectl -n ingress-nginx get svc ingress-nginx-controller -w # 夯住
#    sudo minikube tunnel # 新终端夯住
#    # 第一个命令的EXTERNAL-IP配hosts
#    curl -v http://127.0.0.1:18080 -H 'Host: cozeloop.mini.local' # 验证



chart-tpl-alone-%:
	helm template $(KB_RELEASE_NAME) $(KB_UMBRELLA_PATH) \
		--namespace $(KB_NAMESPACE) \
		-f $(KB_UMBRELLA_PATH)/examples/alone.values.yaml \
		| yq eval '. | select(.kind == "Deployment" and .metadata.name == "coze-loop-$*")' -

chart-tpl-default-%:
	helm template $(KB_RELEASE_NAME) $(KB_UMBRELLA_PATH) \
		--namespace $(KB_NAMESPACE) \
		-f $(KB_UMBRELLA_PATH)/examples/default.values.yaml \
		| yq eval '. | select(.kind == "Deployment" and .metadata.name == "coze-loop-$*")' -

chart-up-%:
	helm upgrade \
    	--install --force $(KB_RELEASE_NAME)-$* $(KB_UMBRELLA_PATH) \
        --namespace $(KB_NAMESPACE) --create-namespace \
        -f $(KB_UMBRELLA_PATH)/values.yaml \
        -f $(KB_UMBRELLA_PATH)/examples/$*.values.yaml

chart:
	@echo "🔧 Building Helm dependencies for umbrella chart..."
	helm dependency build $(KB_UMBRELLA_PATH)

mini-ingress-on:
	minikube addons enable ingress

svc:
	kubectl get svc -n $(KB_DEPLOY_NAME) -o wide

ingress:
	kubectl get ingress -n coze-loop

kb-ctx:
	kubectl config get-contexts
kb-use-ctx-%:
ifeq ($*,mini)
	kubectl config use-context minikube
else
	kubectl config use-context $$(kubectl config get-contexts -o name | grep -v '^minikube$$')
endif
	kubectl config current-context

kb-ns:
	kubectl get namespaces

kb-pod:
	kubectl get pods -n $(KB_NAMESPACE)

kb-del-%:
	helm uninstall $(KB_DEPLOY_NAME)-$* -n $(KB_NAMESPACE)

kb-log-%:
	@echo "Getting logs for the latest pod of $(KB_DEPLOY_NAME)-$* ..."
	@POD=$$(kubectl get pod -n $(KB_NAMESPACE) -l app=$(KB_DEPLOY_NAME)-$* -o jsonpath='{.items[0].metadata.name}'); \
	for c in $$(kubectl get pod $$POD -n $(KB_NAMESPACE) -o jsonpath='{.spec.initContainers[*].name} {.spec.containers[*].name}'); do \
		echo "========== logs from container: $$c =========="; \
		kubectl logs -n $(KB_NAMESPACE) -f $$POD -c $$c --tail=100; \
	done

kb-up-%:
	helm upgrade \
      --install --force $(KB_RELEASE_NAME)-$* $(KB_CHART_PATH)/$* \
      --namespace $(KB_NAMESPACE) --create-namespace \
      -f $(KB_CHART_PATH)/$*/values.yaml && \
    kubectl rollout status deployment/$(KB_DEPLOY_NAME)-$* -n $(KB_NAMESPACE) && \
    POD=$$(kubectl get pod -n $(KB_NAMESPACE) -l app=$(KB_DEPLOY_NAME)-$* -o jsonpath='{.items[0].metadata.name}') && \
    for c in $$(kubectl get pod $$POD -n $(KB_NAMESPACE) -o jsonpath='{.spec.initContainers[*].name} {.spec.containers[*].name}'); do \
      echo "========== logs from container: $$c =========="; \
      kubectl logs -n $(KB_NAMESPACE) -f $$POD -c $$c; \
    done

kb-clean:
	helm list -n $(KB_NAMESPACE) -q | xargs -r -n1 helm uninstall -n $(KB_NAMESPACE)

compose%:
	@case "$*" in \
	  -up) \
	    docker compose \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose.yml \
	      --env-file $(DOCKER_COMPOSE_DIR)/.env \
	      --profile "*" \
	      up ;; \
	  -down) \
	    docker compose \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose.yml \
	      --env-file $(DOCKER_COMPOSE_DIR)/.env \
	      --profile "*" \
	      down ;; \
	  -down-v) \
	    docker compose \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose.yml \
	      --env-file $(DOCKER_COMPOSE_DIR)/.env \
	      --profile "*" \
	      down -v ;; \
	  -up-dev) \
	    docker compose \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose.yml \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose-dev.yml \
	      --env-file $(DOCKER_COMPOSE_DIR)/.env \
	      --profile "*" \
	      up --build ;; \
	  -down-dev) \
	    docker compose \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose.yml \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose-dev.yml \
	      --env-file $(DOCKER_COMPOSE_DIR)/.env \
	      --profile "*" \
	      down ;; \
	  -down-v-dev) \
	    docker compose \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose.yml \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose-dev.yml \
	      --env-file $(DOCKER_COMPOSE_DIR)/.env \
	      --profile "*" \
	      down -v ;; \
	  -up-debug) \
	    docker compose \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose.yml \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose-debug.yml \
	      --env-file $(DOCKER_COMPOSE_DIR)/.env \
	      --profile "*" \
	      up --build ;; \
	  -down-debug) \
	    docker compose \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose.yml \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose-debug.yml \
	      --env-file $(DOCKER_COMPOSE_DIR)/.env \
	      --profile "*" \
	      down ;; \
	  -down-v-debug) \
	    docker compose \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose.yml \
	      -f $(DOCKER_COMPOSE_DIR)/docker-compose-debug.yml \
	      --env-file $(DOCKER_COMPOSE_DIR)/.env \
	      --profile "*" \
	      down -v ;; \
	  -help|*) \
	    echo "Usage:"; \
	    echo "  make compose-up               # up base"; \
	    echo "  make compose-down             # down base"; \
	    echo "  make compose-down-v           # down base + volumes"; \
	    echo "  make compose-up-dev           # up base + dev (build)"; \
	    echo "  make compose-down-dev         # down base + dev"; \
	    echo "  make compose-down-v-dev       # down base + dev + volumes"; \
	    echo "  make compose-up-debug         # up base + debug (build)"; \
	    echo "  make compose-down-debug       # down base + debug"; \
	    echo "  make compose-down-v-debug     # down base + debug + volumes"; \
	    echo; \
	    echo "Notes:"; \
	    echo "  - '--profile \"*\"' is only meaningful for 'up'; it's not required for 'down'."; \
	    echo "  - When you used multiple -f files for 'up', run 'down' with the same -f set."; \
	    exit 1 ;; \
	esac

debug:
	docker compose \
		-f ./release/deployment/docker-compose/docker-compose.yml \
		-f ./release/deployment/docker-compose/debug/docker-compose.yml \
		--env-file ./release/deployment/docker-compose/.env \
		--profile "*" \
		up --build

debug-app:
	docker compose \
		-f ./release/deployment/docker-compose/docker-compose.yml \
		-f ./release/deployment/docker-compose/debug/remote/docker-compose.yml \
		--env-file ./release/deployment/docker-compose/.env \
		--profile "app" \
		up

debug-compose:
	docker compose \
    		-f ./release/deployment/docker-compose/docker-compose.yml \
    		-f ./release/deployment/docker-compose/debug/remote/docker-compose.yml \
    		--profile "*" \
    		config

debug-down-v:
	docker compose \
    		-f ./release/deployment/docker-compose/docker-compose.yml \
    		-f ./release/deployment/docker-compose/debug/docker-compose.yml \
    		--profile "*" \
    		down -v

up:
	docker compose -f ./release/deployment/docker-compose/docker-compose.yml --env-file ./release/deployment/docker-compose/.env --profile "*" up

up-redis:
	docker compose -f ./release/deployment/docker-compose/docker-compose.yml --env-file ./release/deployment/docker-compose/.env --profile "redis" up

up-mysql:
	docker compose -f ./release/deployment/docker-compose/docker-compose.yml --env-file ./release/deployment/docker-compose/.env --profile "mysql" up

up-clickhouse:
	docker compose -f ./release/deployment/docker-compose/docker-compose.yml --env-file ./release/deployment/docker-compose/.env --profile "clickhouse" up

up-minio:
	docker compose -f ./release/deployment/docker-compose/docker-compose.yml --env-file ./release/deployment/docker-compose/.env --profile "minio" up

up-rmq:
	docker compose -f ./release/deployment/docker-compose/docker-compose.yml --env-file ./release/deployment/docker-compose/.env --profile "rmq" up

up-nginx:
	docker compose -f ./release/deployment/docker-compose/docker-compose.yml --env-file ./release/deployment/docker-compose/.env --profile "nginx" up

down:
	docker compose -f ./release/deployment/docker-compose/docker-compose.yml ---profile '*' down

down-v:
	docker compose -f ./release/deployment/docker-compose/docker-compose.yml --profile '*' down -v

image:
	@echo "Building and pushing multi-arch coze-loop images (amd64 + arm64)..."

	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--progress=plain \
		--push \
		-f ./release/image/Dockerfile \
		-t compose-cn-beijing.cr.volces.com/coze/coze-loop:latest \
		.

	@echo "Validating image size from coze-loop:latest (amd64 only)..."

	docker pull compose-cn-beijing.cr.volces.com/coze/coze-loop:latest

	docker run --rm coze-loop:latest du -sh /coze-loop/bin
	docker run --rm coze-loop:latest du -sh /coze-loop/resources
	docker run --rm coze-loop:latest du -sh /coze-loop/conf
	docker run --rm coze-loop:latest du -sh /coze-loop

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
