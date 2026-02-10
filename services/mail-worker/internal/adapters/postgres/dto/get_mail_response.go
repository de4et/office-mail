package dto

type GetMailResponse struct {
	ID   int    `db:"id"`
	To   string `db:"to_addr"`
	From string `db:"from_addr"`
	Body string `db:"body"`
}
