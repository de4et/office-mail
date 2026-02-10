package domain

type Address string

type Mail struct {
	ID   int
	To   Address
	From Address
	Body string
}
