#!/bin/sh

# Ensure database is up before running init script
until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"; do
  >&2 echo "Waiting for PostgreSQL to be available..."
  sleep 5
done

./init-db.sh

# Execute the default command (usually from the Dockerfile)
exec "$@"