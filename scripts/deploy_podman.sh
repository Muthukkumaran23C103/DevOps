#!/usr/bin/env bash
# deploy_podman.sh - Builds container images locally and transfers to target host via ssh | zstd | podman load

set -euo pipefail

TARGET_HOST="${1:-r3-host}"
IMAGE_NAME="${2:-localhost/musicdb-server:latest}"

echo "[+] Building image ${IMAGE_NAME}..."
podman build -t "${IMAGE_NAME}" .

echo "[+] Streaming image to ${TARGET_HOST}..."
podman save "${IMAGE_NAME}" | zstd -T0 -3 | ssh "${TARGET_HOST}" "zstd -d | podman load"

echo "[+] Deployment stream complete on ${TARGET_HOST}."
