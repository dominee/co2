.DEFAULT_GOAL=build

BINARY_NAME=co2


clean:
	go clean
	rm ${BINARY_NAME}
.PHONY:clean

test:
	go test ./...
.PHONY:test

test_coverage:
	go test ./... -coverprofile=coverage.out
.PHONY:test_coverage

dep:
	go mod download
.PHONY:dep

vet: fmt
	go vet
.PHONY:vet

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

build: vet
	go build -o ${BINARY_NAME} co2.go
.PHONY:build
