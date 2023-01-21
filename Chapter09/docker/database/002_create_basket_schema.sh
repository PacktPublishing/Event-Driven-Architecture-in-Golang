#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  CREATE SCHEMA baskets;

  CREATE TABLE baskets.stores_cache
  (
      id         text NOT NULL,
      name       text NOT NULL,
      created_at timestamptz NOT NULL DEFAULT NOW(),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE TRIGGER created_at_stores_trgr BEFORE UPDATE ON baskets.stores_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_stores_trgr BEFORE UPDATE ON baskets.stores_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE baskets.products_cache
  (
      id         text NOT NULL,
      store_id   text NOT NULL,
      name       text NOT NULL,
      price      decimal(9,4) NOT NULL,
      created_at timestamptz NOT NULL DEFAULT NOW(),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      PRIMARY KEY (id)
  );

  CREATE TRIGGER created_at_products_trgr BEFORE UPDATE ON baskets.products_cache FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_products_trgr BEFORE UPDATE ON baskets.products_cache FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE baskets.events
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

  CREATE TABLE baskets.snapshots
  (
      stream_id        text        NOT NULL,
      stream_name      text        NOT NULL,
      stream_version   int         NOT NULL,
      snapshot_name    text        NOT NULL,
      snapshot_data    bytea       NOT NULL,
      updated_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (stream_id, stream_name)
  );

  CREATE TRIGGER updated_at_snapshots_trgr BEFORE UPDATE ON baskets.snapshots FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE baskets.inbox
  (
    id          text NOT NULL,
    name        text NOT NULL,
    subject     text NOT NULL,
    data        bytea NOT NULL,
    received_at timestamptz NOT NULL,
    PRIMARY KEY (id)
  );

  CREATE TABLE baskets.outbox
  (
    id           text NOT NULL,
    name         text NOT NULL,
    subject      text NOT NULL,
    data         bytea NOT NULL,
    published_at timestamptz,
    PRIMARY KEY (id)
  );

  CREATE INDEX basket_unpublished_idx ON baskets.outbox (published_at) WHERE published_at IS NULL;

  GRANT USAGE ON SCHEMA baskets TO mallbots_user;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA baskets TO mallbots_user;
EOSQL
