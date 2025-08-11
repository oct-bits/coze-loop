DOCKER_COMPOSE_DIR := ./release/deployment/docker-compose

HELM_CHART_DIR := ./release/deployment/helm-chart/umbrella
HELM_NAMESPACE := coze-loop
HELM_RELEASE := coze-loop

.PHONY: image mini-start mini-tunnel

.PHONY: FORCE
FORCE:

image:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--progress=plain \
		--push \
		-f ./release/image/Dockerfile \
		-t compose-cn-beijing.cr.volces.com/coze/coze-loop:latest \
		.

	docker pull compose-cn-beijing.cr.volces.com/coze/coze-loop:latest

	docker run --rm coze-loop:latest du -sh /coze-loop/bin
	docker run --rm coze-loop:latest du -sh /coze-loop/resources
	docker run --rm coze-loop:latest du -sh /coze-loop/conf
	docker run --rm coze-loop:latest du -sh /coze-loop

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

helm%:
	@case "$*" in \
	  -chart) \
	    helm dependency build $(HELM_CHART_DIR) ;; \
	  -chart-clean) \
		rm -rf $(HELM_CHART_DIR)/charts $(HELM_CHART_DIR)/Chart.lock ;; \
	  -ctx) \
	    kubectl config get-contexts ;; \
	  -ctx-*) \
		ctx="$*"; \
		ctx="$${ctx#-ctx-}"; \
		echo "switch to context: $$ctx"; \
		kubectl config use-context "$$ctx" ;; \
	  -ns) \
	    kubectl get namespaces ;; \
	  -pod) \
	    kubectl get pods -n $(HELM_NAMESPACE) ;; \
	  -svc) \
	    kubectl get svc -n $(HELM_NAMESPACE) -o wide ;; \
	  -ingress) \
	    kubectl get ingress -n $(HELM_NAMESPACE) ;; \
	  -up) \
		helm upgrade \
		  --install --force $(HELM_RELEASE) $(HELM_CHART_DIR) \
		  --namespace $(HELM_NAMESPACE) --create-namespace \
		  -f $(HELM_CHART_DIR)/values.yaml ;; \
	  -down) \
	    helm list -n $(HELM_NAMESPACE) -q \
	    | \
	    xargs -r -n1 helm uninstall -n $(HELM_NAMESPACE) ;; \
	  -logf-*) \
      	app="$*"; \
      	app="$${app#-logf-}"; \
      	kubectl -n $(HELM_NAMESPACE) logs \
      	  -l app=$(HELM_RELEASE)-$$app \
      	  --all-containers=true \
      	  --tail=100 \
      	  --prefix=true \
		  --max-log-requests=10 \
      	  -f ;; \
	  -tpl-*) \
      	app="$*"; \
      	app="$${app#-tpl-}"; \
      	helm template $(HELM_RELEASE) $(HELM_CHART_DIR) \
      	  --namespace $(HELM_NAMESPACE) \
      	  -f $(HELM_CHART_DIR)/values.yaml | \
      	APP="$$app" yq eval '. | select(.kind == "Deployment" and .metadata.name == ("coze-loop-" + strenv(APP)))' - ;; \
	  --help|*) \
       	echo "Usage:"; \
       	echo "  make helm-chart           # build chart dependencies (helm dependency build)"; \
       	echo "  make helm-chart-clean     # remove chart dependencies"; \
       	echo "  make helm-ctx             # list all kubectl contexts"; \
       	echo "  make helm-ctx-<context>   # switch to a specific kubectl context"; \
       	echo "  make helm-ns              # list all namespaces"; \
       	echo "  make helm-pod             # list all pods in namespace $(HELM_NAMESPACE)"; \
       	echo "  make helm-svc             # list all services in namespace $(HELM_NAMESPACE)"; \
       	echo "  make helm-ingress         # list all ingress resources in namespace $(HELM_NAMESPACE)"; \
       	echo "  make helm-up              # upgrade/install release $(HELM_RELEASE) from chart"; \
       	echo "  make helm-down            # uninstall all releases in namespace $(HELM_NAMESPACE)"; \
       	echo "  make helm-logf-<app>      # follow logs of all containers in pods with app=$(HELM_RELEASE)-<app>"; \
       	echo "  make helm-tpl-<app>       # render Deployment manifest of coze-loop-<app> locally"; \
       	echo; \
       	echo "Notes:"; \
       	echo "  - Ensure $(HELM_NAMESPACE) and $(HELM_RELEASE) are set before running commands."; \
       	echo "  - Commands with '-<name>' suffix accept a dynamic argument (e.g., helm-ctx-xxx, helm-logf-app)."; \
       	echo "  - '-tpl-*' renders manifests without applying them to the cluster."; \
       	exit 1 ;; \
	esac

mini-start:
	minikube start --addons=ingress

mini-tunnel:
	minikube tunnel
