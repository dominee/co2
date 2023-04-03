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

test_coverage:
	@echo "$$(tput bold)Coverage test$$(tput sgr0)"
	go test ./... -coverprofile=coverage.out
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

run: build
	@echo "$$(tput bold)Running...$$(tput sgr0)"
	cat testdata/test.txt | ./${BINARY_NAME}
.PHONY:run