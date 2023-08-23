package sqx

import "fmt"

// ErrTooManyRows indicates that a query returned more rows than expected. This is used in calls to OneStrict() which
// expects a single row to be returned. In Strict mode, this error is raised if the number of rows returned is not equal
// to the expected number. If you received this error in your code and didn't expect it, check out the One() or First()
// methods instead.
type ErrTooManyRows struct {
	Expected int
	Actual   int
}

func (e ErrTooManyRows) Error() string {
	return fmt.Errorf("too many rows: expected = %d actual = %d",
		e.Expected, e.Actual).Error()
}
