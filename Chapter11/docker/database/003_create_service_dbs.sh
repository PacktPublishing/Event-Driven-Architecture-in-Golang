#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE commondb TEMPLATE template0;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "commondb" <<-EOSQL
  -- Apply to keep modifications to the created_at column from being made
  CREATE OR REPLACE FUNCTION created_at_trigger()
  RETURNS TRIGGER AS \$\$
  BEGIN
    NEW.created_at := OLD.created_at;
    RETURN NEW;
  END;
  \$\$ language 'plpgsql';

  -- Apply to a table to automatically update update_at columns
  CREATE OR REPLACE FUNCTION updated_at_trigger()
  RETURNS TRIGGER AS \$\$
  BEGIN
     IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
        NEW.updated_at = NOW();
        RETURN NEW;
     ELSE
        RETURN OLD;
     END IF;
  END;
  \$\$ language 'plpgsql';
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE baskets TEMPLATE commondb;

  CREATE USER baskets_user WITH ENCRYPTED PASSWORD 'baskets_pass';
  GRANT USAGE ON SCHEMA public TO baskets_user;
  GRANT CREATE, CONNECT ON DATABASE baskets TO baskets_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "baskets" <<-EOSQL
  CREATE SCHEMA baskets;
  GRANT CREATE, USAGE ON SCHEMA baskets TO baskets_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE cosec TEMPLATE commondb;

  CREATE USER cosec_user WITH ENCRYPTED PASSWORD 'cosec_pass';
  GRANT USAGE ON SCHEMA public TO cosec_user;
  GRANT CREATE, CONNECT ON DATABASE cosec TO cosec_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "cosec" <<-EOSQL
  CREATE SCHEMA cosec;
  GRANT CREATE, USAGE ON SCHEMA cosec TO cosec_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE customers TEMPLATE commondb;

  CREATE USER customers_user WITH ENCRYPTED PASSWORD 'customers_pass';
  GRANT USAGE ON SCHEMA public TO customers_user;
  GRANT CREATE, CONNECT ON DATABASE customers TO customers_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "customers" <<-EOSQL
  CREATE SCHEMA customers;
  GRANT CREATE, USAGE ON SCHEMA customers TO customers_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE depot TEMPLATE commondb;

  CREATE USER depot_user WITH ENCRYPTED PASSWORD 'depot_pass';
  GRANT USAGE ON SCHEMA public TO depot_user;
  GRANT CREATE, CONNECT ON DATABASE depot TO depot_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "depot" <<-EOSQL
  CREATE SCHEMA depot;
  GRANT CREATE, USAGE ON SCHEMA depot TO depot_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE notifications TEMPLATE commondb;

  CREATE USER notifications_user WITH ENCRYPTED PASSWORD 'notifications_pass';
  GRANT USAGE ON SCHEMA public TO notifications_user;
  GRANT CREATE, CONNECT ON DATABASE notifications TO notifications_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "notifications" <<-EOSQL
  CREATE SCHEMA notifications;
  GRANT CREATE, USAGE ON SCHEMA notifications TO notifications_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE ordering TEMPLATE commondb;

  CREATE USER ordering_user WITH ENCRYPTED PASSWORD 'ordering_pass';
  GRANT USAGE ON SCHEMA public TO ordering_user;
  GRANT CREATE, CONNECT ON DATABASE ordering TO ordering_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "ordering" <<-EOSQL
  CREATE SCHEMA ordering;
  GRANT CREATE, USAGE ON SCHEMA ordering TO ordering_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE payments TEMPLATE commondb;

  CREATE USER payments_user WITH ENCRYPTED PASSWORD 'payments_pass';
  GRANT USAGE ON SCHEMA public TO payments_user;
  GRANT CREATE, CONNECT ON DATABASE payments TO payments_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "payments" <<-EOSQL
  CREATE SCHEMA payments;
  GRANT CREATE, USAGE ON SCHEMA payments TO payments_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE search TEMPLATE commondb;

  CREATE USER search_user WITH ENCRYPTED PASSWORD 'search_pass';
  GRANT USAGE ON SCHEMA public TO search_user;
  GRANT CREATE, CONNECT ON DATABASE search TO search_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "search" <<-EOSQL
  CREATE SCHEMA search;
  GRANT CREATE, USAGE ON SCHEMA search TO search_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE stores TEMPLATE commondb;

  CREATE USER stores_user WITH ENCRYPTED PASSWORD 'stores_pass';
  GRANT USAGE ON SCHEMA public TO stores_user;
  GRANT CREATE, CONNECT ON DATABASE stores TO stores_user;
EOSQL
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "stores" <<-EOSQL
  CREATE SCHEMA stores;
  GRANT CREATE, USAGE ON SCHEMA stores TO stores_user;
EOSQL
