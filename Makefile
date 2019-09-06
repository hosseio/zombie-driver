.PHONY: all test

all:
    docker-compose up -d
	make -C ./driver-location
	make -C ./gateway
	make -C ./zombie-driver

test:
    docker-compose up -d
	make -C ./driver-location test
	make -C ./gateway test
	make -C ./zombie-driver test
