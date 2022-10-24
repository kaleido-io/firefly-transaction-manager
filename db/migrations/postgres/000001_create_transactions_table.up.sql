CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

BEGIN;
CREATE TABLE transactions (
  seq                      INTEGER         PRIMARY KEY AUTOINCREMENT,
  id                       VARCHAR(64)     NOT NULL,
  tx_status                TEXT            NOT NULL,
  sequence_id              UUID            NOT NULL,
  nonce                    BIGINT          NOT NULL,
  gas                      BIGINT          NOT NULL,
  transaction_headers      TEXT            NOT NULL,
  transaction_data         VARCHAR(64)     NOT NULL,
  transaction_hash         VARCHAR(64),
  gas_price                TEXT            NOT NULL, 
  policy_info              TEXT            NOT NULL, 
  first_submit             BIGINT,
  last_submit              BIGINT,
  receipt                  TEXT,
  error_msg                TEXT,
  error_history            TEXT            NOT NULL,
  confirmations            TEXT,
  updated                  BIGINT,
  created                  BIGINT          NOT NULL,
);

CREATE UNIQUE INDEX transactions_id on transactions(id);
-- TODO: create index for nonce + signer?
COMMIT;