#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  INSERT INTO customers.customers (id, name, sms_number, enabled)
    VALUES ('f0e2d41a-a485-4008-b578-747732ae1089', 'Buyer #1', '555-1212', true);
EOSQL
