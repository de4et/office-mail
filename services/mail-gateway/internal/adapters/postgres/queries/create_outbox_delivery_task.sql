INSERT into outbox (aggregate_id, created_at, status, retry_at, trace_context)
VALUES ($1, NOW(), 'new', NOW(), $2);
