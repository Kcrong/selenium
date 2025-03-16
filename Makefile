.PHONY: all test lint format

# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

test:
	$(GOTEST) -v -race -cover ./...

lint:
	golangci-lint run --config .golangci.yml

format:
	@go install -v github.com/incu6us/goimports-reviser/v3@latest
	golangci-lint run --fix --config .golangci.yml
	@goimports-reviser -rm-unused \
		-company-prefixes 'github.com/Kcrong,github.com/hodlgap' \
		-project-name 'github.com/Kcrong/selenigo' \
		-format \
		-set-alias \
		./...
	@gofmt -s -w .
	$(GOMOD) tidy