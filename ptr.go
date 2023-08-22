package sqx

// Ptr is a convenience method for converting inline constants into pointers for use with ToClause and ToSetMap
func Ptr[T any](t T) *T {
	return &t
}
