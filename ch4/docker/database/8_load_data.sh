#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  INSERT INTO stores.stores (id, name, location, participating)
    VALUES ('bd9e26cd-d861-4b74-862f-06e0825c7af3', 'Waldorf Books', 'West Upper Level', false);

  INSERT INTO stores.products (id, store_id, name, description, sku, price)
    VALUES ('557e2f05-b10a-475c-baff-2cea8a86d8df', 'bd9e26cd-d861-4b74-862f-06e0825c7af3', 'EDA with Golang', '', '1234', 49.99);

  INSERT INTO customers.customers (id, name, sms_number, enabled)
    VALUES ('f0e2d41a-a485-4008-b578-747732ae1089', 'Buyer #1', '555-1212', true);

EOSQL
