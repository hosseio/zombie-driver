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
	@echo "    init                 - start the driver-location container and its dependencies"
	@echo "    finish               - stop the driver-location container and its dependencies"
	@echo "    dependencies         - start the containers which the driver-location needs"
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

dependencies:
	cd ".."
	docker-compose up -d redis
	docker-compose up -d nsqlookupd
	docker-compose up -d nsqd

init: dependencies up

finish:
	cd ".."
	docker-compose down

up:
	cd ".."
	docker-compose up -d driver-location

down:
	cd ".."
	docker-compose down driver-location