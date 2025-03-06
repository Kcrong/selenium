vet:
	go vet ./...

.PHONY: format
format:
	@go install golang.org/x/tools/cmd/goimports@latest
	goimports -local "github.com/Kcrong/selenium" -w .
	gofmt -s -w .
	go mod tidy
