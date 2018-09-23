BINARY=bin/schools_api

USERID=$(shell id -u)
GROUPID=$(shell id -g)

.PHONY: all vendor build run

all: vendor build

vendor:
	docker-compose run govendor dep ensure
	docker-compose run govendor chown -R $(USERID):$(GROUPID) .

build:
	docker-compose run gobuilder go build -o $(BINARY)
	docker-compose run gobuilder chown -R $(USERID):$(GROUPID) $(BINARY)

run: build
	docker-compose up schoolsapi
