INSERT into outbox (aggregate_id, created_at, status, retry_at)
VALUES ($1, NOW(), 'new', NOW());
