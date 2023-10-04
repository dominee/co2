.DEFAULT_GOAL=build

BINARY_NAME=co2


clean:
	@echo "$$(tput bold)Cleaning$$(tput sgr0)"
	go clean
.PHONY:clean

test:
	@echo "$$(tput bold)Testing$$(tput sgr0)"
	go test ./...
.PHONY:test

test-coverage:
	@echo "$$(tput bold)Coverage test$$(tput sgr0)"
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
.PHONY:test_coverage

dep:
	@echo "$$(tput bold)Downloading dependencies$$(tput sgr0)"
	go mod download
.PHONY:dep

vet: fmt
	@echo "$$(tput bold)Vetting$$(tput sgr0)"
	go vet
.PHONY:vet

fmt:
	@echo "$$(tput bold)Formating$$(tput sgr0)"
	go fmt ./...
.PHONY:fmt

lint: fmt
	@echo "$$(tput bold)Linting$$(tput sgr0)"
	golint ./...
.PHONY:lint

build: vet
	@echo "$$(tput bold)Building$$(tput sgr0)"
	go build -o ${BINARY_NAME} co2.go
.PHONY:build

build-all: vet
	@echo "$$(tput bold)Building all targets$$(tput sgr0)"
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin ${BINARY_NAME}.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux ${BINARY_NAME}.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows.exe ${BINARY_NAME}.go
.PHONY:build-all

run: build
	@echo "$$(tput bold)Running...$$(tput sgr0)"
	cat testdata/test.txt | ./${BINARY_NAME}
.PHONY:run