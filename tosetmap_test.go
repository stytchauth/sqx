package sqx_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stytchauth/sqx"
)

func ptr[T any](v T) *T { return &v }

type thingyUpdateFilter struct {
	StrCol *string `db:"str_col"`
	IntCol *int    `db:"int_col"`
}

func ExampleToSetMap() {
	type filter struct {
		Value  *string   `db:"first_col"`
		Values *[]string `db:"second_col"`
	}
	setMap, _ := sqx.ToSetMap(&filter{
		Value:  ptr("example"),
		Values: &[]string{"a", "b"},
	})

	out, _ := json.Marshal(setMap)
	fmt.Print(string(out))

	// Output:
	// {"first_col":"example","second_col":["a","b"]}
}

func TestToSetMap(t *testing.T) {
	filter := thingyUpdateFilter{
		StrCol: ptr("i am str"),
		IntCol: ptr(100),
	}
	t.Run("Can convert all fields of a struct to a map", func(t *testing.T) {
		expected := map[string]any{
			"str_col": ptr("i am str"),
			"int_col": ptr(100),
		}

		setMap, err := sqx.ToSetMap(&filter)
		assert.NoError(t, err)
		assert.Equal(t, expected, setMap)
	})

	filter2 := thingyUpdateFilter{
		StrCol: ptr("still a str"),
	}
	t.Run("Omits specific fields when asked", func(t *testing.T) {
		expected := map[string]any{
			"str_col": ptr("still a str"),
		}

		setMap, err := sqx.ToSetMap(&filter2, "str_ptr_col_null", "int_col")
		assert.NoError(t, err)
		assert.Equal(t, expected, setMap)
	})

	t.Run("Returns empty on nil input", func(t *testing.T) {
		setMap, err := sqx.ToSetMap(nil)
		assert.Equal(t, setMap, map[string]any{})
		assert.NoError(t, err)
	})
	t.Run("Returns empty on empty input", func(t *testing.T) {
		emptyUpdate := &thingyUpdateFilter{}
		setMap, err := sqx.ToSetMap(emptyUpdate)
		assert.Equal(t, setMap, map[string]any{})
		assert.NoError(t, err)
	})
	t.Run("Returns empty on struct without any DB tags", func(t *testing.T) {
		emptyUpdate := &struct{}{}
		setMap, err := sqx.ToSetMap(emptyUpdate)
		assert.Equal(t, setMap, map[string]any{})
		assert.NoError(t, err)
	})
}

func TestToSetMapAlias(t *testing.T) {
	f1 := thingyUpdateFilter{
		StrCol: ptr("i am str"),
		IntCol: ptr(100),
	}
	expected := map[string]any{
		"table.str_col": ptr("i am str"),
		"table.int_col": ptr(100),
	}

	setMap, err := sqx.ToSetMapAlias("table", &f1)
	assert.NoError(t, err)
	assert.Equal(t, expected, setMap)
}
