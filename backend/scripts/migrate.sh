#!/bin/bash
set -e

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-clearn}
DB_PASSWORD=${DB_PASSWORD:-clearn_secret}
DB_NAME=${DB_NAME:-clearn}

export PGPASSWORD="$DB_PASSWORD"

echo "=== Running Auth migrations ==="
for f in migrations/auth/*.up.sql; do
    echo "  Applying $f..."
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$f"
done

echo "=== Running LMS migrations ==="
for f in migrations/lms/*.up.sql; do
    echo "  Applying $f..."
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$f"
done

echo "=== Running Gamification migrations ==="
for f in migrations/gamification/*.up.sql; do
    echo "  Applying $f..."
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$f"
done

echo "=== All migrations applied successfully ==="
