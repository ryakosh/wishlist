#!bin/sh

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER wishlist;
	CREATE DATABASE wishlist;
	GRANT ALL PRIVILEGES ON DATABASE wishlist TO wishlist;
EOSQL