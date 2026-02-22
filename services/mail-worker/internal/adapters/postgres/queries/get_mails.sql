SELECT id, body, from_addr, to_addr
FROM mail
WHERE id = ANY($1)
