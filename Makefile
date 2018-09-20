BINARY=bin/schools_api

build:
	go build -o $(BINARY)

.PHONY: run
run: bin/schools_api
	$(BINARY)
