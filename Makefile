.PHONY: all test

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
	@echo "    e2e-test             - rebuild the whole app and run and e2e test"
	@echo "    up                   - start the driver-location container"
	@echo "    down                 - stop the driver-location container"
	@echo

test: unit-test integration-test

unit-test:
	make -C ./driver-location unit-test
	make -C ./gateway unit-test
	make -C ./zombie-driver unit-test

e2e-test: up
	go mod download
	go mod verify
	go test -p=1 -count=1 -v ./... -tags=e2e

integration-test:
	make -C ./driver-location integration-test
	make -C ./gateway integration-test
	make -C ./zombie-driver integration-test

install:
	make -C ./driver-location install
	make -C ./gateway install
	make -C ./zombie-driver install

up:
	docker-compose up -d --force-recreate --build

down:
	docker-compose down