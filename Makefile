APP_VERSION := $(shell grep 'appVersion:' deploy/charts/weaviate-health/Chart.yaml | awk '{print $$2}')

.PHONY: compose_up
compose_up: test build_image
	@docker-compose up

.PHONY: build_image
build_image:
	@echo "Building Docker Image with tag: $(or ${tag},$(APP_VERSION))"
	@docker build . -t adamkisala/health:$(or ${tag},$(APP_VERSION)) --file deploy/weaviate-health.dockerfile --no-cache
	@docker push adamkisala/health:$(or ${tag},$(APP_VERSION))

.PHONY: test
test: ## Run all tests
	$(call blue, "# Running tests...")
	go test ./...
	helm template deploy/charts/weaviate-health --values deploy/charts/weaviate-health/values.yaml
