SHELL := /usr/bin/bash
.ONESHELL:

build: main.go $(wildcard **/*.go)
	TAG=$$(git tag | tail -1)
	COMMIT=$$(git rev-parse HEAD)
	VERSION=$${TAG}.$${COMMIT: -8}
	TIMESTAMP=$$(date +%s)
	OUT=api.il2missionplanner.com.v$${VERSION}.$${TIMESTAMP}.out
	go build -o $${OUT} -ldflags "-X main.version=$${VERSION}" main.go
	shasum $${OUT}

develop:
	go run main.go

test:
	go test . ./config ./handlers ./server -cover -v

clean:
	rm *.out
