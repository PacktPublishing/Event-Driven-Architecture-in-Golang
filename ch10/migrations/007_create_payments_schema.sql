-- +goose Up
CREATE SCHEMA payments;

SET SEARCH_PATH TO payments, public;

CREATE TABLE payments (
  id          text          NOT NULL,
  customer_id text          NOT NULL,
  amount      decimal(9, 4) NOT NULL,
  created_at  timestamptz   NOT NULL DEFAULT NOW(),
  updated_at  timestamptz   NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_payments_trgr
  BEFORE UPDATE
  ON payments
  FOR EACH ROW
EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_payments_trgr
  BEFORE UPDATE
  ON payments
  FOR EACH ROW
EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE invoices (
  id         text          NOT NULL,
  order_id   text          NOT NULL,
  amount     decimal(9, 4) NOT NULL,
  status     text          NOT NULL,
  created_at timestamptz   NOT NULL DEFAULT NOW(),
  updated_at timestamptz   NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE INDEX invoices_order_id_idx ON invoices (order_id);

CREATE TRIGGER created_at_invoices_trgr
  BEFORE UPDATE
  ON invoices
  FOR EACH ROW
EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_invoices_trgr
  BEFORE UPDATE
  ON invoices
  FOR EACH ROW
EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE inbox (
  id          text        NOT NULL,
  name        text        NOT NULL,
  subject     text        NOT NULL,
  data        bytea       NOT NULL,
  received_at timestamptz NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE outbox (
  id           text  NOT NULL,
  name         text  NOT NULL,
  subject      text  NOT NULL,
  data         bytea NOT NULL,
  published_at timestamptz,
  PRIMARY KEY (id)
);

CREATE INDEX payments_unpublished_idx ON outbox (published_at) WHERE published_at IS NULL;

-- +goose Down
DROP SCHEMA IF EXISTS payments CASCADE;
-- SET SEARCH_PATH TO payments;
--
-- DROP TABLE IF EXISTS outbox;
-- DROP TABLE IF EXISTS inbox;
-- DROP TABLE IF EXISTS invoices;
-- DROP TABLE IF EXISTS payments;
