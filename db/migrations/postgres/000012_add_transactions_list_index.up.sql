BEGIN;
-- Composite index to optimize queries filtering by from address, status, and nonce
-- This index significantly improves fetching unprocessed transactions for a specific from address with nonce filtering
CREATE INDEX transactions_from_status_nonce ON transactions(tx_from, status, tx_nonce);

-- Index on created timestamp to optimize queries filtering or sorting by creation time
CREATE INDEX transactions_created ON transactions(created);

COMMIT;

