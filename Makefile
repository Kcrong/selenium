.PHONY: format
format:
	@#go install golang.org/x/tools/cmd/goimports@latest
	@#goimports -local "github.com/hodlgap/captive-portal" -w .
	@go install -v github.com/incu6us/goimports-reviser/v3@latest
	@goimports-reviser -rm-unused \
		-company-prefixes 'github.com/Kcrong' \
		-excludes 'db' \
		-project-name 'github.com/Kcrong/selenium' \
		-format \
		./...
	@gofmt -s -w .
	@go mod tidy

.PHONY: lint
lint:
	@golangci-lint run -v ./...
