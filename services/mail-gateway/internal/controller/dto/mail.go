package dto

type Mail struct {
	ID   int    `json:"id,omitempty"`
	To   string `json:"to"`
	From string `json:"from"`
	Body string `json:"body"`
}
