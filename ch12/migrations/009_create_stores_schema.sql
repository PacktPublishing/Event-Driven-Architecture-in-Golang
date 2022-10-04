-- +goose Up
CREATE SCHEMA stores;

SET
SEARCH_PATH TO stores, PUBLIC;

CREATE TABLE stores (
  id            text        NOT NULL,
  name          text        NOT NULL,
  location      text        NOT NULL,
  participating bool        NOT NULL DEFAULT FALSE,
  created_at    timestamptz NOT NULL DEFAULT NOW(),
  updated_at    timestamptz NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE INDEX participating_stores_idx ON stores (participating) WHERE participating;

CREATE TRIGGER created_at_stores_trgr
  BEFORE UPDATE
  ON stores
  FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_stores_trgr
  BEFORE UPDATE
  ON stores
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE products (
  id          text          NOT NULL,
  store_id    text          NOT NULL,
  name        text          NOT NULL,
  description text          NOT NULL,
  sku         text          NOT NULL,
  price       decimal(9, 4) NOT NULL,
  created_at  timestamptz   NOT NULL DEFAULT NOW(),
  updated_at  timestamptz   NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE INDEX store_products_idx ON products (store_id);

CREATE TRIGGER created_at_products_trgr
  BEFORE UPDATE
  ON products
  FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_products_trgr
  BEFORE UPDATE
  ON products
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE events (
  stream_id      text        NOT NULL,
  stream_name    text        NOT NULL,
  stream_version int         NOT NULL,
  event_id       text        NOT NULL,
  event_name     text        NOT NULL,
  event_data     bytea       NOT NULL,
  occurred_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (stream_id, stream_name, stream_version)
);

CREATE TABLE snapshots (
  stream_id      text        NOT NULL,
  stream_name    text        NOT NULL,
  stream_version int         NOT NULL,
  snapshot_name  text        NOT NULL,
  snapshot_data  bytea       NOT NULL,
  updated_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (stream_id, stream_name)
);

CREATE TRIGGER updated_at_snapshots_trgr
  BEFORE UPDATE
  ON snapshots
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE inbox (
  id          text        NOT NULL,
  name        text        NOT NULL,
  subject     text        NOT NULL,
  data        bytea       NOT NULL,
  metadata    bytea       NOT NULL,
  sent_at     timestamptz NOT NULL,
  received_at timestamptz NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE outbox (
  id           text        NOT NULL,
  name         text        NOT NULL,
  subject      text        NOT NULL,
  data         bytea       NOT NULL,
  metadata     bytea       NOT NULL,
  sent_at      timestamptz NOT NULL,
  published_at timestamptz,
  PRIMARY KEY (id)
);

CREATE INDEX stores_unpublished_idx ON outbox (published_at) WHERE published_at IS NULL;

-- +goose Down
DROP SCHEMA IF EXISTS stores CASCADE;
-- SET SEARCH_PATH TO stores;
--
-- DROP TABLE IF EXISTS outbox;
-- DROP TABLE IF EXISTS inbox;
-- DROP TABLE IF EXISTS snapshots;
-- DROP TABLE IF EXISTS events;
-- DROP TABLE IF EXISTS products;
-- DROP TABLE IF EXISTS stores;
