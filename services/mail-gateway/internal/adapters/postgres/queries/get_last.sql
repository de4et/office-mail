SELECT
  id,
  body,
  from_addr,
  to_addr
FROM
  mail
WHERE from_addr=$1
ORDER BY id ASC
LIMIT $2 OFFSET $3
