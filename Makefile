.ONESHELL:

SHELL := /bin/bash
.SHELLFLAGS := -ec

format: ## Fixes format errors by applying the canonical Go style
	gofmt -d -s -w pkg

test:
	docker-compose up -d
	go test -gcflags '-N -l' -cover -v -race -p 1 ./...
	docker logs wiremock-golang-client_wiremock_1
	docker-compose down

generate: clean
	go generate ./...

clean:
	find $(CURDIR) -name "*_easyjson.go" -delete