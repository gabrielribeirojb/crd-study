#!/usr/bin/env bash
set -euo pipefail

MODULE="github.com/gabrielribeirojb/crd-study"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

CODEGEN_PKG="${ROOT_DIR}/vendor/k8s.io/code-generator"

bash "${CODEGEN_PKG}/kube_codegen.sh" \
  --output-base "${ROOT_DIR}" \
  --go-header-file "${ROOT_DIR}/hack/boilerplate.go.txt" \
  --boilerplate "${ROOT_DIR}/hack/boilerplate.go.txt" \
  --modules "${MODULE}" \
  --apis-dir "${MODULE}/pkg/apis" \
  --groups "gdch:v1" \
  --generate "deepcopy,client,informer,lister"