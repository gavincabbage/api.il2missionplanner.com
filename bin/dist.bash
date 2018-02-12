#!/usr/bin/env bash


set -x

git fetch --tags

tag=$(git tag | tail -1)
commit=$(git rev-parse HEAD)
version=${tag}.${commit: -8}
timestamp=$(date +%s)
appname=api.il2missionplanner.com
outfile=${appname}_v${version}.${timestamp}.out

mkdir -p dist
go build -o ./dist/${outfile} -ldflags "-X main.version=${version}" main.go
cd ./dist/ # go into dir to keep hash file free of "./dist/"
shasum ${outfile} > ${outfile}.sha256

exit 0