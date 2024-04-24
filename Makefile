APP_VERSION := $(shell grep 'appVersion:' deploy/charts/health/Chart.yaml | awk '{print $$2}')

.PHONY: compose_up
compose_up: test build_test_image
	@docker-compose up

.PHONY: build_test_image
build_test_image:
	@echo "Building test Docker Image with tag: latest"
	@docker build . -t health:latest --file deploy/health.dockerfile --no-cache

.PHONY: build_image
build_image:
	@echo "Building Docker Image with tag: $(or ${tag},$(APP_VERSION))"
	@docker build . -t adamkisala/health:$(or ${tag},$(APP_VERSION)) --file deploy/health.dockerfile --no-cache
	@docker push adamkisala/health:$(or ${tag},$(APP_VERSION))

.PHONY: test
test: ## Run all tests
	$(call blue, "# Running tests...")
	go test ./...
	helm template deploy/charts/health --values deploy/charts/health/values.yaml
