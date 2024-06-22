echo "Create Database keycloak"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$POSTGRES_DB" <<-EOSQL
    DROP DATABASE IF EXISTS keycloak;
	CREATE DATABASE keycloak with encoding 'UTF8';
EOSQL
