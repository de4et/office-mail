package postgres

import "github.com/lib/pq"

func IsDuplicate(err *pq.Error) bool {
	return err.Code == "23505"
}
