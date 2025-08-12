HELM_CHART_DIR := .
HELM_NAMESPACE := coze-loop
HELM_RELEASE := coze-loop

.PHONY: FORCE
FORCE:

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

chart-%:
	@helm package $(HELM_CHART_DIR) --version $*
	@helm push $(HELM_RELEASE)-$*.tgz oci://compose-cn-beijing.cr.volces.com/coze
	@rm -f $(HELM_RELEASE)-$*.tgz