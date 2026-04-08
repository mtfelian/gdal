#!/usr/bin/env bash
set -euo pipefail

image_name="${IMAGE_NAME:-gdal-tests}"
cache_prefix="${CACHE_PREFIX:-gdal-tests}"

docker build -f Dockerfile.tests -t "${image_name}" .

MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker run --rm \
  --mount "type=volume,src=${cache_prefix}-go-build-cache,dst=/root/.cache/go-build" \
  --mount "type=volume,src=${cache_prefix}-go-mod-cache,dst=/go/pkg/mod" \
  "${image_name}" go test -v -count=1 ./... "$@"
