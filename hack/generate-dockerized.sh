#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

docker run --rm -v $PWD:/workspace -w /workspace -e GOPROXY=https://goproxy.cn,direct golang:1.15 hack/generate.sh
