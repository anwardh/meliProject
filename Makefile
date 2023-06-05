PACKAGES_PATH = $(shell go list -f '{{ .Dir }}' ./...)

.PHONY: start
start:
	@go run cmd/server/main.go

