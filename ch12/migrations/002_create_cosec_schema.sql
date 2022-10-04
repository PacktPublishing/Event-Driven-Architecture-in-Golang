-- +goose Up
CREATE SCHEMA cosec;

SET
SEARCH_PATH TO cosec, PUBLIC;

CREATE TABLE sagas (
  id           text        NOT NULL,
  name         text        NOT NULL,
  data         bytea       NOT NULL,
  step         int         NOT NULL,
  done         bool        NOT NULL,
  compensating bool        NOT NULL,
  updated_at   timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id, name)
);

CREATE TRIGGER updated_at_co_sagas_trgr
  BEFORE UPDATE
  ON sagas
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

CREATE INDEX cosec_unpublished_idx ON baskets.outbox (published_at) WHERE published_at IS NULL;

-- +goose Down
DROP SCHEMA IF EXISTS cosec CASCADE;
-- SET SEARCH_PATH TO cosec;
--
-- DROP TABLE IF EXISTS outbox;
-- DROP TABLE IF EXISTS inbox;
-- DROP TABLE IF EXISTS sagas;
