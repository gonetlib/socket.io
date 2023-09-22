.PHONY: all
all:
	go install ./...

.PHONY: test
test:
	go clean -testcache && go test -race -cover -covermode=atomic ./...

.PHONY: bench
bench:
	go clean -testcache && go test -bench . -benchmem ./...

.PHONY: lint
lint:
	golangci-lint run 

.PHONY: cover
cover:
	go clean -testcache && go test ./... -cover -coverprofile=c.out && go tool cover -html=c.out