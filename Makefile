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
