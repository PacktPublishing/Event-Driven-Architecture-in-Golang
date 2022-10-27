-- +goose Up
CREATE SCHEMA depot;

SET
SEARCH_PATH TO depot, PUBLIC;

CREATE TABLE stores_cache (
  id         text        NOT NULL,
  name       text        NOT NULL,
  location   text        NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_stores_trgr
  BEFORE UPDATE
  ON stores_cache
  FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_stores_trgr
  BEFORE UPDATE
  ON stores_cache
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE products_cache (
  id         text        NOT NULL,
  store_id   text        NOT NULL,
  name       text        NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_products_trgr
  BEFORE UPDATE
  ON products_cache
  FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_products_trgr
  BEFORE UPDATE
  ON products_cache
  FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE shopping_lists (
  id              text        NOT NULL,
  order_id        text        NOT NULL,
  stops           bytea       NOT NULL,
  assigned_bot_id text        NOT NULL,
  status          text        NOT NULL,
  created_at      timestamptz NOT NULL DEFAULT NOW(),
  updated_at      timestamptz NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE INDEX shopping_lists_order_id_idx ON shopping_lists (order_id);
CREATE INDEX shopping_lists_availability_idx ON shopping_lists (status, created_at) WHERE status = 'available';

CREATE TRIGGER created_at_shopping_lists_trgr
  BEFORE UPDATE
  ON shopping_lists
  FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_shopping_lists_trgr
  BEFORE UPDATE
  ON shopping_lists
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

CREATE INDEX depot_unpublished_idx ON outbox (published_at) WHERE published_at IS NULL;

-- +goose Down
DROP SCHEMA IF EXISTS depot CASCADE;
-- SET SEARCH_PATH TO depot;
--
-- DROP TABLE IF EXISTS outbox;
-- DROP TABLE IF EXISTS inbox;
-- DROP TABLE IF EXISTS shopping_lists;
-- DROP TABLE IF EXISTS products_cache;
-- DROP TABLE IF EXISTS stores_cache;
