#!/bin/ash

set -e
echo "run db migrations"

/app/migrate -path ./migration -database "$DB_SOURCE" -verbose up

echo "start app"
exec "$@"
