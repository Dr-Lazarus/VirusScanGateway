# Makefile
.PHONY: start

start:
	@echo "Starting the server..."
	@go run cmd/server/main.go
