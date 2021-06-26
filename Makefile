.PHONY: all prepare test lint build

prepare:
	@go mod download

build:
	@go build -o euro21

run:
	@go run main.go