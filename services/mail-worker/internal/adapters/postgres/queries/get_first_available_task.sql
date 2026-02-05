select
  id,
  aggregate_id
from
  outbox
where
  status = 'new'
limit
  1;
