BEGIN;
-- Composite index to optimize queries filtering by from address, and nonce for pending transactions
-- This index significantly improves fetching unprocessed transactions for a specific from address with nonce filtering
CREATE INDEX transactions_from_pending_nonce ON transactions(tx_from, tx_nonce) WHERE status = 'Pending';

-- Index on created timestamp to optimize queries filtering or sorting by creation time
CREATE INDEX transactions_created ON transactions(created);

COMMIT;

