#!/usr/bin/bash
set -xe

export GRAFANA_VERSION="10.4.0pre"

yarn install --immutable
yarn build
go run build.go setup
go run build.go build
nfpm pkg --packager deb --config ./nfpm.yaml \
  --target ./dist/grafana-x86_64-amd64.deb
