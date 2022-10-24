BEGIN;

CREATE TABLE streams (
  seq                      INTEGER         PRIMARY KEY AUTOINCREMENT,
  id                       UUID            NOT NULL,
  stream_name              VARCHAR(64),
  suspended                BOOLEAN,
  stream_type              VARCHAR(64),     
  error_handling           VARCHAR(64)     NOT NULL,
  -- is bigint fine for uint64?
  batch_size               BIGINT          NOT NULL,           
  batch_timeout            BIGINT          NOT NULL,           
  retry_timeout            BIGINT          NOT NULL,           
  blocked_retry_delay      BIGINT          NOT NULL,           
  eth_batch_timeout        BIGINT, 
  eth_retry_timeout        BIGINT,
  eth_blocked_retry_delay  BIGINT,
  webhook                  TEXT,
  websocket                TEXT,
  updated                  BIGINT,
  created                  BIGINT          NOT NULL,
);

CREATE UNIQUE INDEX streams_id on streams(id);
COMMIT;