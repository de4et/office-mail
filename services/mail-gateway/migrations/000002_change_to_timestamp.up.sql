ALTER TABLE outbox
  ADD COLUMN created_at_ts TIMESTAMPTZ,
  ADD COLUMN retry_at_ts TIMESTAMPTZ;

UPDATE outbox
SET
  created_at_ts = created_at::timestamptz,
  retry_at_ts   = retry_at::timestamptz;

ALTER TABLE outbox
  ALTER COLUMN created_at_ts SET DEFAULT now();

ALTER TABLE outbox
  DROP COLUMN created_at,
  DROP COLUMN retry_at;

ALTER TABLE outbox
  RENAME COLUMN created_at_ts TO created_at;

ALTER TABLE outbox
  RENAME COLUMN retry_at_ts TO retry_at;

