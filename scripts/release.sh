#!/bin/sh

curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
./bin/goreleaser release

