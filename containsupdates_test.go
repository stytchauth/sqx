package sqx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stytchauth/sqx"
)

func ExampleContainsUpdates() {
	type filter struct {
		Value  *string   `db:"first_col"`
		Values *[]string `db:"second_col"`
	}
	first := sqx.ContainsUpdates(&filter{
		Value:  ptr("example"),
		Values: &[]string{"a", "b"},
	})
	second := sqx.ContainsUpdates(&filter{ /* Empty! */ })
	fmt.Printf("first: %t second: %t", first, second)

	// Output:
	// first: true second: false
}

func TestContainsUpdates(t *testing.T) {
	type fooUpdateFilter struct {
		StrCol *string  `db:"str_col"`
		IntCol *int     `db:"int_col"`
		PtrCol **string `db:"ptr_col"`
	}

	t.Run("Returns true when fields are non-nil", func(t *testing.T) {
		filter := fooUpdateFilter{
			StrCol: ptr("i am str"),
			IntCol: ptr(1),
		}
		assert.True(t, sqx.ContainsUpdates(&filter))
	})

	t.Run("Returns true when pointer to pointer is set, but nil", func(t *testing.T) {
		var nilString *string = nil
		filter := fooUpdateFilter{
			PtrCol: ptr(nilString),
		}
		assert.True(t, sqx.ContainsUpdates(&filter))
	})

	t.Run("Omits unset fields", func(t *testing.T) {
		emptyFilter := fooUpdateFilter{}
		assert.False(t, sqx.ContainsUpdates(&emptyFilter))
	})
}
