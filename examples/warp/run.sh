#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -eq 0 ]; then
  set -- /app/testdata/tiles.gpkg /tmp/tiles-3857.gpkg -t_srs epsg:3857 -of GPKG
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec "${script_dir}/../run-example.sh" warp "$@"
