.PHONY: image_test
image_test: build_image_test

.PHONY: compose_up ## Downloads dependencies, runs tests, builds image based on {TAG}, and runs docker-compose
compose_up: test image_test
	@docker-compose up

# Testing targets
.PHONY: build_image_test
build_image_test:
	@echo "Building Docker Image test"
	@docker build . -t weaviate-health:test --file deploy/weaviate-health.dockerfile --no-cache

.PHONY: test
test: ## Run all tests
	$(call blue, "# Running tests...")
	go test ./...
	helm template deploy/charts/weaviate-health --values deploy/charts/weaviate-health/values.yaml
