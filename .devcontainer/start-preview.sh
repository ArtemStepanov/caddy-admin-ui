#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"

for _ in $(seq 1 60); do
  if docker info >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

if ! docker info >/dev/null 2>&1; then
  echo "Docker did not become ready within 60 seconds" >&2
  exit 1
fi

args=(up --detach --wait --wait-timeout 300)
if [[ "${1:-}" == "--build" ]]; then
  args+=(--build)
fi

docker compose -f compose.preview.yml "${args[@]}"
docker compose -f compose.preview.yml ps
