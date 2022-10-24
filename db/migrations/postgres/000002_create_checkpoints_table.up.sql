BEGIN;

CREATE TABLE checkpoints (
  seq         INTEGER         PRIMARY KEY AUTOINCREMENT,
  streamId    UUID            NOT NULL,
  listeners   TEXT            NOT NULL 
  created     BIGINT          NOT NULL,
);

CREATE UNIQUE INDEX checkpoints_streamId on checkpoints(streamId);
COMMIT;