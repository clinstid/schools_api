BINARY=bin/schools_api

USERID=$(shell id -u)
GROUPID=$(shell id -g)

.PHONY: all vendor build run

all: vendor build test down

vendor:
	docker-compose run --rm govendor dep ensure
	docker-compose run --rm govendor chown -R $(USERID):$(GROUPID) .

build:
	docker-compose run --rm gobuilder go build -o $(BINARY)
	docker-compose run --rm gobuilder chown -R $(USERID):$(GROUPID) bin

run: build
	docker-compose up schoolsapi

down:
	docker-compose down

test:
	docker-compose run --rm tester
	docker-compose run --rm tester chown -R $(USERID):$(GROUPID) .

clean: down
	docker-compose run --rm tester \
	    find . \
		-name __pycache__ -o \
		-name .pytest_cache -o \
		-name "*.pyc" \
		-exec rm -rvf {} \;
