RUN_CONTEXT ?= docker compose run --rm go

shell:
	$(RUN_CONTEXT) bash

run:
	$(RUN_CONTEXT) go run main.go

fmt:
	$(RUN_CONTEXT) go fmt ./...

lint:
	$(RUN_CONTEXT) golangci-lint run

mod/tidy:
	$(RUN_CONTEXT) go mod tidy
