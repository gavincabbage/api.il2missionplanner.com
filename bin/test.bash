#!/usr/bin/env bash

repodir="api.il2missionplanner.com"
wd=$(basename $(pwd))
if [ "${wd}" != "${repodir}" ]
then
    echo "error! run tests from ${repodir} directory"
    exit 1
else
    # Any package with tests will need to be added to this list
    go test . ./config ./handlers ./server
fi

