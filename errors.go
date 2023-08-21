package sqx

import "fmt"

type ErrTooManyRows struct {
	Expected int
	Actual   int
}

func (e *ErrTooManyRows) Error() string {
	return fmt.Errorf("too many rows: expected = %d actual = %d",
		e.Expected, e.Actual).Error()
}
