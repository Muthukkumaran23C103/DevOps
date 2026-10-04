#!/usr/bin/env bash
# backup_app_db.sh - Nightly automated backup script for user app data (ratings, reviews, lists)
# Isolates the 'app' schema to protect user data from catalog dumps/resets.

set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-/var/backups/musicdb}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/musicdb_app_schema_${TIMESTAMP}.sql.gz"
DB_CONTAINER="${DB_CONTAINER:-musicdb-postgres}"
DB_USER="${DB_USER:-musicdb_app}"
DB_NAME="${DB_NAME:-musicdb}"

mkdir -p "${BACKUP_DIR}"

echo "[+] Starting dump of 'app' schema at ${TIMESTAMP}..."

podman exec "${DB_CONTAINER}" pg_dump -U "${DB_USER}" -d "${DB_NAME}" --schema=app --format=plain --clean --if-exists | gzip > "${BACKUP_FILE}"

echo "[+] App schema backup completed successfully: ${BACKUP_FILE}"

# Keep last 30 daily backups
find "${BACKUP_DIR}" -name "musicdb_app_schema_*.sql.gz" -mtime +30 -delete
echo "[+] Cleaned up backups older than 30 days."
