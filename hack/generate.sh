#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

go install sigs.k8s.io/controller-tools/cmd/controller-gen

controller-gen object crd output:crd:dir=./deploy/crd paths=./operator/...
