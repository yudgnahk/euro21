.PHONY: all prepare test lint build

prepare:
	@go mod download

table:
	@go run main.go table

build:
	@go build -o euro21