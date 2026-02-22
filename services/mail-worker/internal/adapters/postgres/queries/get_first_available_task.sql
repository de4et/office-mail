select
  id,
  aggregate_id,
  trace_context
from
  outbox
where
  status = 'new'
limit
  $1
FOR UPDATE SKIP LOCKED;
