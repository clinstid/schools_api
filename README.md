# schools_api

[![Build Status](https://travis-ci.org/clinstid/schools_api.svg?branch=master)](https://travis-ci.org/clinstid/schools_api)

API for Managing a List of Schools

# API Details

This is an example HTTP API implemented in Go with the [Gin framework](https://github.com/gin-gonic/gin) that supports a few operations on a list of schools (colleges in the United States):
- `GET /schools`: Retrieves a paginated list of schools with `name` and `id`
- `POST /schools`: Adds a new school to the list
- `GET /schools/:schoolId`: Retrieve the school with the specified `schoolId`
- `PUT /schools/:schoolId`: Updates the school with the specified `schoolId`

The full API specification is available at [./spec/schools_api_spec.yaml](./spec/schools_api_spec.yaml).

# Developing, Testing, and Running the API Service

The tools for running the API require GNU make, [docker](https://www.docker.com/get-started), and [docker-compose](https://docs.docker.com/compose/install/) and assume you are either running on Linux or macOS.

To build and test the API:
```sh
make all
```

That runs the following make targets:
- `vendor`: Uses godep to install the Go dependencies
- `build`: Builds the `schools_api` binary
- `test`: Starts the HTTP service and runs python functional tests against it
- `down`: Tears down the docker-compose services

If you want to run the service and play around:
```sh
make run
```

The HTTP service will then be available at [http://localhost:8080/schools](http://localhost:8080/schools). Interactive documentation (using [swagger-ui](https://swagger.io/tools/swagger-ui/)) will be available at [http://localhost:8080/docs](http://localhost:8080/docs).
