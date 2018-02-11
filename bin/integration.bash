#!/usr/bin/env bash


set -x

go run -ldflags "-X main.version=integration" main.go &
sleep 2
lsof_out=$(lsof -i :9999 -Fp | head -1)
api_pid=${lsof_out: 1}
echo "API PID: ${api_pid}"

newman run \
        test/api.il2missionplanner.com.postman_collection.json \
        -e test/api.il2missionplanner.com_local.postman_environment.json

kill -9 ${api_pid}