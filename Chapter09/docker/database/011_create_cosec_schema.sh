#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  CREATE SCHEMA cosec;

  CREATE TABLE cosec.sagas
  (
      id           text        NOT NULL,
      name         text        NOT NULL,
      data         bytea       NOT NULL,
      step         int         NOT NULL,
      done         bool        NOT NULL,
      compensating bool        NOT NULL,
      updated_at   timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (id, name)
  );

  CREATE TRIGGER updated_at_co_sagas_trgr BEFORE UPDATE ON cosec.sagas FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE cosec.inbox
  (
    id          text NOT NULL,
    name        text NOT NULL,
    subject     text NOT NULL,
    data        bytea NOT NULL,
    received_at timestamptz NOT NULL,
    PRIMARY KEY (id)
  );

  CREATE TABLE cosec.outbox
  (
    id           text NOT NULL,
    name         text NOT NULL,
    subject      text NOT NULL,
    data         bytea NOT NULL,
    published_at timestamptz,
    PRIMARY KEY (id)
  );

  CREATE INDEX cosec_unpublished_idx ON baskets.outbox (published_at) WHERE published_at IS NULL;

  GRANT USAGE ON SCHEMA cosec TO mallbots_user;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA cosec TO mallbots_user;
EOSQL
