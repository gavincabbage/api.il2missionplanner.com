#!/usr/bin/env bash


set -x

git fetch --tags

TAG=$(git tag | tail -1)
COMMIT=$(git rev-parse HEAD)
VERSION=${TAG}.${COMMIT: -8}
TIMESTAMP=$(date +%s)
OUT=api.il2missionplanner.com.v${VERSION}.${TIMESTAMP}.out

mkdir -p dist
go build -o ./dist/${OUT} -ldflags "-X main.version=${VERSION}" main.go
shasum ./dist/${OUT} > ./dist/${OUT}.sha256

exit 0