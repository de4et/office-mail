update outbox set status='done'
where id=ANY($1);
