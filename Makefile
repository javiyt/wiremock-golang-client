.ONESHELL:

SHELL := /bin/bash
.SHELLFLAGS := -ec
.DEFAULT_GOAL := test

format: ## Fixes format errors by applying the canonical Go style
	gofmt -d -s -w pkg

test: generate
	docker-compose up -d
	go test -gcflags '-N -l' -cover -v -race -p 1 ./...
	docker-compose down

generate: clean
	go generate ./...

clean:
	find $(CURDIR) -name "*_easyjson.go" -delete