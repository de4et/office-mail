INSERT INTO mail (body, from_addr, to_addr)
VALUES ($1, $2, $3)
RETURNING id;
