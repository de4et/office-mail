package domain

type OutboxTask struct {
	ID          int
	AggregateID int
}
