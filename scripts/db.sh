#!/usr/bin/env bash
set -euo pipefail

PG_HOST=${PG_HOST:-}
PG_USER=${PG_USER:-}
PG_PASSWORD=${PG_PASSWORD:-}
PG_PORT=${PG_PORT:-5432}

if [ -z "${PG_HOST}" ] || [ -z "${PG_USER}" ]; then
	echo "Could not connect to database"
  echo "Database host or user cannot be empty"
  exit 1
fi

function usage() {
    echo "Usage:"
    echo "  $0 create DATABASE"
    exit 1
}

test -z "${1-}" && usage
command="$1"
shift

test -z "${1-}" && usage
database="$1"
shift

create=$(cat <<EOF
    DO
    \$body\$
    BEGIN
      IF NOT EXISTS (
        SELECT * FROM pg_catalog.pg_user WHERE usename = 'rami'
      ) THEN
        CREATE USER rami WITH PASSWORD 'rami';
      END IF;
    END
    \$body\$;

    CREATE DATABASE $database ENCODING 'UTF-8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8' TEMPLATE template0 OWNER rami;
    \c $database

    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO rami;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO rami;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO rami;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO rami;
EOF
)

drop=$(cat <<EOF
    DROP DATABASE IF EXISTS $database;
EOF
)

terminate=$(cat <<EOF
    SELECT pg_terminate_backend(pg_stat_activity.pid)
    FROM pg_stat_activity
    WHERE pg_stat_activity.datname = '$database' AND pid <> pg_backend_pid();
EOF
)

case "$command" in
  "create")
    echo "$create" | PGPASSWORD=${PG_PASSWORD} psql -h${PG_HOST} -U${PG_USER} -p${PG_PORT} -v ON_ERROR_STOP=1
    ;;
  "drop")
    echo "$terminate" "$drop" | PGPASSWORD=${PG_PASSWORD} psql -h${PG_HOST} -U${PG_USER} -p${PG_PORT} -v ON_ERROR_STOP=1
    ;;
esac

exit $?