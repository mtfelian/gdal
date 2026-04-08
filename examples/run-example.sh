#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -lt 1 ]; then
  echo "usage: run-example.sh <example-dir> [args...]" >&2
  exit 1
fi

example_dir="$1"
shift

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"

image_name="${IMAGE_NAME:-gdal-tests}"
cache_prefix="${CACHE_PREFIX:-gdal-tests}"

(
  cd "${repo_root}"
  docker build -f Dockerfile.tests -t "${image_name}" .
)

MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker run --rm \
  --mount "type=volume,src=${cache_prefix}-go-build-cache,dst=/root/.cache/go-build" \
  --mount "type=volume,src=${cache_prefix}-go-mod-cache,dst=/go/pkg/mod" \
  -w "/app/examples/${example_dir}" \
  "${image_name}" go run . "$@"
