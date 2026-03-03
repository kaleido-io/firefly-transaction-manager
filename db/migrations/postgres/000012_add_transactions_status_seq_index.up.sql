BEGIN;
CREATE INDEX IF EXISTS transactions_status_seq ON transactions(status, seq);
COMMIT;
