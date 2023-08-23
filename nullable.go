package sqx

type Nullable[T any] **T

// NewNullable creates a Nullable[T] from a provided value
// use it to set nullable fields in Update calls to a concrete value
func NewNullable[T any](t T) Nullable[T] {
	return Ptr(Ptr(t))
}

// NewNull creates a Nullable[T] from a provided value
// use it to set nullable fields in Update calls to a null value
func NewNull[T any]() Nullable[T] {
	return Ptr[*T](nil)
}
