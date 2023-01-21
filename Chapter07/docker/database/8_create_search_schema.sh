#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  CREATE SCHEMA search;

  CREATE TABLE search.customers_cache
  (
      id         text NOT NULL,
      name       text NOT NULL,
      created_at timestamptz NOT NULL DEFAULT NOW(),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE TRIGGER created_at_customers_trgr BEFORE UPDATE ON search.customers_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_customers_trgr BEFORE UPDATE ON search.customers_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE search.stores_cache
  (
      id         text NOT NULL,
      name       text NOT NULL,
      created_at timestamptz NOT NULL DEFAULT NOW(),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE TRIGGER created_at_stores_trgr BEFORE UPDATE ON search.stores_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_stores_trgr BEFORE UPDATE ON search.stores_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE search.products_cache
  (
      id         text NOT NULL,
      store_id   text NOT NULL,
      name       text NOT NULL,
      created_at timestamptz NOT NULL DEFAULT NOW(),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE TRIGGER created_at_products_trgr BEFORE UPDATE ON search.products_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_products_trgr BEFORE UPDATE ON search.products_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE search.orders
  (
    order_id       text NOT NULL,
    customer_id    text NOT NULL,
    customer_name  text NOT NULL,
    items          bytea NOT NULL,
    status         text NOT NULL,
    product_ids    text ARRAY NOT NULL,
    store_ids      text ARRAY NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT NOW(),
    updated_at     timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (order_id)
  );

  -- CREATE TRIGGER created_at_orders_trgr BEFORE UPDATE ON ordering.orders FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_orders_trgr BEFORE UPDATE ON ordering.orders FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();


  GRANT USAGE ON SCHEMA search TO mallbots_user;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA search TO mallbots_user;
EOSQL
