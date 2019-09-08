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
	@echo "    up                   - start the driver-location container"
	@echo "    down                 - stop the driver-location container"
	@echo

test: unit-test integration-test

unit-test:
	make -C ./driver-location unit-test
	make -C ./gateway unit-test
	make -C ./zombie-driver unit-test

integration-test:
	make -C ./driver-location integration-test
	make -C ./gateway integration-test
	make -C ./zombie-driver integration-test

install:
	make -C ./driver-location install
	make -C ./gateway install
	make -C ./zombie-driver install

up:
	docker-compose up -d --force-recreate

down:
	docker-compose down --rmi all --remove-orphans