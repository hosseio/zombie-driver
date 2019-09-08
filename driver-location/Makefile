.PHONY: all help install unit-test integration-test test init dependencies up down

all: help

help:
	@echo
	@echo "usage: make <command>"
	@echo
	@echo "commands:"
	@echo "    install              - populate vendor/ from go.mod"
	@echo "    unit-test            - run unit tests"
	@echo "    integration-test     - run integration tests"
	@echo "    test                 - run unit test and integration tests"
	@echo "    up                   - start the driver-location container"
	@echo "    down                 - stop the driver-location container"
	@echo

test: unit-test integration-test

unit-test:
	go test -p=1 -count=1 -v ./...

integration-test:
	go test -p=1 -count=1 -v ./... -tags=integration

install:
	@go mod vendor
	@go mod tidy

up:
	cd ".."
	@docker-compose up -d --force-recreate driver-location

down:
	cd ".."
	@docker-compose stop driver-location