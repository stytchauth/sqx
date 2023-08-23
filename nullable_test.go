package sqx

import (
	"fmt"
	"reflect"
)

func ExampleNewNullable() {
	type updateFilter struct {
		Field1 Nullable[int] `db:"field_1"`
	}
	sm, _ := ToSetMap(&updateFilter{
		Field1: NewNullable(1),
	})

	fmt.Printf("setting field_1 to %v", reflect.ValueOf(sm["field_1"]).Elem().Elem().Interface())
	// Output:
	// setting field_1 to 1
}

func ExampleNewNull() {
	type updateFilter struct {
		Field2 Nullable[string] `db:"field_2"`
	}
	sm, _ := ToSetMap(&updateFilter{
		Field2: NewNull[string](),
	})

	fmt.Printf("setting field_2 to %v", reflect.ValueOf(sm["field_2"]).Elem().Interface())
	// Output:
	// setting field_2 to <nil>
}
