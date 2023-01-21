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

  CREATE TABLE stores.events
  (
      stream_id      text        NOT NULL,
      stream_name    text        NOT NULL,
      stream_version int         NOT NULL,
      event_id       text        NOT NULL,
      event_name     text        NOT NULL,
      event_data     bytea       NOT NULL,
      occurred_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (stream_id, stream_name, stream_version)
  );

  CREATE TABLE stores.snapshots
  (
      stream_id        text        NOT NULL,
      stream_name      text        NOT NULL,
      stream_version   int         NOT NULL,
      snapshot_name    text        NOT NULL,
      snapshot_data    bytea       NOT NULL,
      updated_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (stream_id, stream_name)
  );

  CREATE TRIGGER updated_at_snapshots_trgr BEFORE UPDATE ON stores.snapshots FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE stores.inbox
  (
    id          text NOT NULL,
    name        text NOT NULL,
    subject     text NOT NULL,
    data        bytea NOT NULL,
    received_at timestamptz NOT NULL,
    PRIMARY KEY (id)
  );

  CREATE TABLE stores.outbox
  (
    id           text NOT NULL,
    name         text NOT NULL,
    subject      text NOT NULL,
    data         bytea NOT NULL,
    published_at timestamptz,
    PRIMARY KEY (id)
  );

  CREATE INDEX stores_unpublished_idx ON stores.outbox (published_at) WHERE published_at IS NULL;

  GRANT USAGE ON SCHEMA stores TO mallbots_user;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA stores TO mallbots_user;
EOSQL
