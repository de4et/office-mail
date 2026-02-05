create table mail (
  id SERIAL PRIMARY KEY,
  body TEXT,
  from_addr TEXT,
  to_addr TEXT
);

create table outbox (
  id SERIAL PRIMARY KEY,
  aggregate_id INT,
  created_at DATE,
  status TEXT,
  retry_at DATE
);
