BEGIN;
CREATE INDEX transactions_status_seq ON transactions(status, seq);
COMMIT;
