#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mallbots" <<-EOSQL
  CREATE SCHEMA payments;

  CREATE TABLE payments.payments
  (
    id          text NOT NULL,
    customer_id text NOT NULL,
    amount      decimal(9, 4) NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
  );

  CREATE TRIGGER created_at_payments_trgr BEFORE UPDATE ON payments.payments FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_payments_trgr BEFORE UPDATE ON payments.payments FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  CREATE TABLE payments.invoices
  (
    id         text NOT NULL,
    order_id   text NOT NULL,
    amount     decimal(9,4) NOT NULL,
    status     text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
  );

  CREATE INDEX invoices_order_id_idx ON payments.invoices (order_id);

  CREATE TRIGGER created_at_invoices_trgr BEFORE UPDATE ON payments.invoices FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
  CREATE TRIGGER updated_at_invoices_trgr BEFORE UPDATE ON payments.invoices FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

  GRANT USAGE ON SCHEMA payments TO mallbots_user;
  GRANT INSERT, UPDATE, DELETE, SELECT ON ALL TABLES IN SCHEMA payments TO mallbots_user;
EOSQL
