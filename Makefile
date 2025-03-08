.PHONY: format
format:
	@go install -v github.com/incu6us/goimports-reviser/v3@latest
	@goimports-reviser -rm-unused \
		-company-prefixes 'github.com/Kcrong' \
		-project-name 'github.com/Kcrong/selenium' \
		-format \
		./...
	@gofmt -s -w .
	@go mod tidy

.PHONY: lint
lint:
	@golangci-lint run -v ./...

.PHONY: integration

# Detect system architecture
ARCH := $(shell uname -m)
DOCKER_COMPOSE_CMD := docker compose

# Select profile based on architecture
ifeq ($(ARCH),arm64)
	DOCKER_COMPOSE_PROFILE := --profile enable-on-arm
else
	DOCKER_COMPOSE_PROFILE := --profile disable-on-arm
endif

# Run integration tests with the correct profile
integration:
	@echo "Running integration tests on architecture: $(ARCH)"
	@$(DOCKER_COMPOSE_CMD) $(DOCKER_COMPOSE_PROFILE) up --build --abort-on-container-exit --exit-code-from integration-test
	@$(DOCKER_COMPOSE_CMD) down --volumes --remove-orphans
