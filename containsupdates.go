package sqx

import (
	scan "github.com/blockloop/scan/v2"
)

// ContainsUpdates returns true if an update filter is nonempty.
// This function panics if v is not a pointer to a struct.
func ContainsUpdates(v any, excluded ...string) bool {
	if isNil(v) {
		return false
	}
	cols, err := scan.ColumnsStrict(v, excluded...)
	if err != nil {
		// Err will only be returned if v is not a pointer to a struct
		// so panics should only ever occur in development (assuming code is ran)
		panic(err)
	}
	vals, err := scan.Values(cols, v)
	if err != nil {
		// Err will only be returned if v is not a pointer to a struct
		// so panics should only ever occur in development (assuming code is ran)
		panic(err)
	}
	for _, val := range vals {
		if !isNil(val) {
			return true
		}
	}
	return false
}
