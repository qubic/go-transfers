# go transfers service

This integration layer service provides asset and qu transfer information to integrators. The information is based
on qubic events (aka 'logging').

## Requirements

* Go (golang)
* Postgresql database

## Setup

### Create database user and schema

```postgresql
CREATE USER some_user_name WITH PASSWORD 'some-password';
CREATE DATABASE some_database OWNER some_user_name;
```

Database migrations use [golang-migrate](https://github.com/golang-migrate/migrate) and are automatically applied on 
startup. For manual migration you can use migrate-cli (see golang-migrate docs), for example to clean the database use
`down` or `drop`:

```shell
migrate -source file://path/to/migrations -database postgres://localhost:5432/database down
```

### Configure environment

Environment variables need to be set or necessary values need to be passed as command line arguments.

## Build & Run

Run `go build` in the root folder. Then you can run the executable.

### Run tests

Run `go test -v ./...` to execute all tests. To exclude system integration tests that have dependencies to external
systems and therefore need configuration use the `ci` tag : `go test -v -tags ci ./...`.