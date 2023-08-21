package sqx

import (
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T { return &v }

func ExampleToClause() {
	type filter struct {
		Value  *string   `db:"first_col"`
		Values *[]string `db:"second_col"`
	}
	clause := ToClause(&filter{
		Value:  ptr("example"),
		Values: &[]string{"a", "b"},
	})
	sql, args, _ := clause.ToSql()
	fmt.Printf("%s, %s", sql, args)

	// Output:
	// first_col = ? AND second_col IN (?,?), [example a b]
}

type thingyGetFilter struct {
	StrCol *string `db:"str_col"`
	IntCol *[]int  `db:"int_col"`
}

func TestToClause(t *testing.T) {
	filter := thingyGetFilter{
		StrCol: ptr("i am str"),
		IntCol: &[]int{1, 2},
	}
	t.Run("Can convert all fields of a struct to a map", func(t *testing.T) {
		expected := sq.Eq{
			"str_col": ptr("i am str"),
			"int_col": &[]int{1, 2},
		}

		clause := ToClause(&filter)
		assert.Equal(t, expected, clause.contents)
	})

	filter2 := thingyGetFilter{
		StrCol: ptr("still a str"),
	}
	t.Run("Omits unset fields", func(t *testing.T) {
		expected := sq.Eq{
			"str_col": ptr("still a str"),
		}

		clause := ToClause(&filter2)
		assert.Equal(t, expected, clause.contents)
	})

	t.Run("Returns empty on nil input", func(t *testing.T) {
		expected := sq.Eq{}

		clause := ToClause(nil)
		assert.Equal(t, expected, clause.contents)
	})
	t.Run("Returns empty on empty output", func(t *testing.T) {
		expected := sq.Eq{}

		clause := ToClause(&thingyGetFilter{})
		assert.Equal(t, expected, clause.contents)
	})
}

func TestToClauseAlias(t *testing.T) {
	f1 := thingyGetFilter{
		StrCol: ptr("i am str"),
		IntCol: &[]int{100},
	}
	expected := sq.Eq{
		"table.str_col": ptr("i am str"),
		"table.int_col": &[]int{100},
	}

	clause := ToClauseAlias("table", &f1)
	assert.Equal(t, expected, clause.contents)
}
