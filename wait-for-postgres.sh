#!/bin/sh

# Usage: ./wait-for-postgres.sh db ./main

set -e

host="$1"
shift
cmd="$@"

echo "Waiting for PostgreSQL at $host..."

until nc -z "$host" 5432; do
  sleep 1
done

echo "PostgreSQL is up. Starting app..."
exec $cmd
