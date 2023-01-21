#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  CREATE SCHEMA stores;

  CREATE TABLE stores.stores
  (
      id            text NOT NULL,
      name          text NOT NULL,
      location      text NOT NULL,
      participating bool NOT NULL DEFAULT FALSE,
      created_at    timestamptz NOT NULL DEFAULT NOW(),
      updated_at    timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE INDEX participating_stores_idx ON stores.stores (participating) WHERE participating;

  CREATE TRIGGER created_at_stores_trgr BEFORE UPDATE ON stores.stores FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_stores_trgr BEFORE UPDATE ON stores.stores FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE stores.products
  (
      id          text NOT NULL,
      store_id    text NOT NULL,
      name        text NOT NULL,
      description text NOT NULL,
      sku         text NOT NULL,
      price       decimal(9,4) NOT NULL,
      created_at  timestamptz NOT NULL DEFAULT NOW(),
      updated_at  timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE INDEX store_products_idx ON stores.products (store_id);

  CREATE TRIGGER created_at_products_trgr BEFORE UPDATE ON stores.products FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_products_trgr BEFORE UPDATE ON stores.products FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  GRANT USAGE ON SCHEMA stores TO mallbots_user;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA stores TO mallbots_user;
EOSQL
