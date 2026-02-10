package dto

type GetFirstAvailableTaskResponse struct {
	ID           int    `db:"id"`
	AggregateID  int    `db:"aggregate_id"`
	TraceContext []byte `db:"trace_context"`
}
