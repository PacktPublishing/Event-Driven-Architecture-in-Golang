#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  CREATE SCHEMA ordering;

  CREATE TABLE ordering.orders
  (
    id          text NOT NULL,
    customer_id text NOT NULL,
    payment_id  text NOT NULL,
    invoice_id  text NOT NULL,
    shopping_id text NOT NULL,
    items       bytea NOT NULL,
    status      text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
  );

  CREATE TRIGGER created_at_orders_trgr BEFORE UPDATE ON ordering.orders FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_orders_trgr BEFORE UPDATE ON ordering.orders FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE ordering.events
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

  CREATE TABLE ordering.snapshots
  (
      stream_id        text        NOT NULL,
      stream_name      text        NOT NULL,
      stream_version   int         NOT NULL,
      snapshot_name    text        NOT NULL,
      snapshot_data    bytea       NOT NULL,
      updated_at       timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (stream_id, stream_name)
  );

  CREATE TRIGGER updated_at_snapshots_trgr BEFORE UPDATE ON ordering.snapshots FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE ordering.inbox
  (
    id          text NOT NULL,
    name        text NOT NULL,
    subject     text NOT NULL,
    data        bytea NOT NULL,
    received_at timestamptz NOT NULL,
    PRIMARY KEY (id)
  );

  CREATE TABLE ordering.outbox
  (
    id           text NOT NULL,
    name         text NOT NULL,
    subject      text NOT NULL,
    data         bytea NOT NULL,
    published_at timestamptz,
    PRIMARY KEY (id)
  );

  CREATE INDEX ordering_unpublished_idx ON ordering.outbox (published_at) WHERE published_at IS NULL;

  GRANT USAGE ON SCHEMA ordering TO mallbots_user;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA ordering TO mallbots_user;
EOSQL
