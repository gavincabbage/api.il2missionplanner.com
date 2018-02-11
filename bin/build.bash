#!/usr/bin/env bash

tag=$(git tag | tail -1)
commit=$(git rev-parse HEAD)
version=${tag}.${commit: -8}
date=$(date +%d%m%Y)
go build -o api.il2missionplanner.com.${version}.${date}.out -ldflags "-X main.version=${version}" *.go