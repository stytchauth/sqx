package sqx

// EmptyResult represents a result with no rows affected.
// This is used for an UpdateBuilder that has no pending changes since the query would be a noop.
type EmptyResult struct{}

func (e EmptyResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (e EmptyResult) RowsAffected() (int64, error) {
	return 0, nil
}
