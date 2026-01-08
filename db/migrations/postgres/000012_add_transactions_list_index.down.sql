BEGIN;
DROP INDEX IF EXISTS transactions_from_pending_nonce;
DROP INDEX IF EXISTS transactions_created;
COMMIT;

