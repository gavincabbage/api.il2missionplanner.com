SHELL := bash
.ONESHELL:

develop:
	go run -ldflags "-X main.version=develop" main.go

test: unit integration

unit:
	go test . ./config ./handlers ./server -cover -v

integration:
	./bin/integration.bash

dist: main.go $(wildcard **/*.go) test
	./bin/dist.bash

clean:
	rm -rf ./dist/
