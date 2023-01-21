#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  CREATE SCHEMA depot;

  CREATE TABLE depot.stores_cache
  (
      id         text NOT NULL,
      name       text NOT NULL,
      location   text NOT NULL,
      created_at timestamptz NOT NULL DEFAULT NOW(),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE TRIGGER created_at_stores_trgr BEFORE UPDATE ON depot.stores_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_stores_trgr BEFORE UPDATE ON depot.stores_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE depot.products_cache
  (
      id         text NOT NULL,
      store_id   text NOT NULL,
      name       text NOT NULL,
      created_at timestamptz NOT NULL DEFAULT NOW(),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE TRIGGER created_at_products_trgr BEFORE UPDATE ON depot.products_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_products_trgr BEFORE UPDATE ON depot.products_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE depot.shopping_lists
  (
      id              text NOT NULL,
      order_id        text NOT NULL,
      stops           bytea NOT NULL,
      assigned_bot_id text NOT NULL,
      status          text NOT NULL,
      created_at      timestamptz NOT NULL DEFAULT NOW(),
      updated_at      timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE INDEX shopping_lists_order_id_idx ON depot.shopping_lists (order_id);
  CREATE INDEX shopping_lists_availability_idx ON depot.shopping_lists (status, created_at) WHERE status = 'available';

  CREATE TRIGGER created_at_shopping_lists_trgr BEFORE UPDATE ON depot.shopping_lists FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_shopping_lists_trgr BEFORE UPDATE ON depot.shopping_lists FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  GRANT USAGE ON SCHEMA depot TO mallbots_user;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA depot TO mallbots_user;
EOSQL
