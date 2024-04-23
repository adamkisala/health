.PHONY: compose_up
compose_up: test build_image
	@docker-compose up

.PHONY: build_image
build_image:
	@echo "Building Docker Image test"
	@docker build . -t adamkisala/health:latest --file deploy/weaviate-health.dockerfile --no-cache
	@docker push adamkisala/health:latest

.PHONY: test
test: ## Run all tests
	$(call blue, "# Running tests...")
	go test ./...
	helm template deploy/charts/weaviate-health --values deploy/charts/weaviate-health/values.yaml

.PHONY: upgrade
upgrade: build_image
	$(call blue, "# Upgrading helm chart...")
	helm upgrade weaviate-health ./deploy/charts/weaviate-health \
      --values ./deploy/charts/weaviate-health/values.yaml \
      --install \
      --create-namespace \
      --namespace monitoring \
      --wait \
      --timeout 2m \
      --reset-values \
      --cleanup-on-fail \
      --atomic \

.PHONY: dry_run_upgrade
dry_run_upgrade: ## Upgrade the helm chart dry-run
	$(call blue, "# Upgrading helm chart dry-run...")
	helm upgrade weaviate-health ./deploy/charts/weaviate-health \
      --values ./deploy/charts/weaviate-health/values.yaml \
      --install \
      --create-namespace \
      --namespace monitoring \
      --wait \
      --timeout 5m \
      --reset-values \
      --cleanup-on-fail \
      --atomic \
      --dry-run
