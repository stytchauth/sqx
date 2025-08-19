package sqx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleToClause() {
	type filter struct {
		Value  *string   `db:"first_col"`
		Values *[]string `db:"second_col"`
	}
	clause := ToClause(&filter{
		Value:  Ptr("example"),
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

type thingyGetFilterWithNoTags struct {
	StrCol *string
	IntCol *[]int
}

type thingyGetFilterWithDuplicateTags struct {
	Field1 *string `db:"same_col"`
	Field2 *string `db:"same_col"`
}

func TestToClause(t *testing.T) {
	filter := thingyGetFilter{
		StrCol: Ptr("i am str"),
		IntCol: &[]int{1, 2},
	}
	t.Run("Can convert all fields of a struct to a map", func(t *testing.T) {
		expected := Eq{
			"str_col": Ptr("i am str"),
			"int_col": &[]int{1, 2},
		}

		clause := ToClause(&filter)
		assert.Equal(t, expected, clause.contents)
	})

	filter2 := thingyGetFilter{
		StrCol: Ptr("still a str"),
	}
	t.Run("Omits unset fields", func(t *testing.T) {
		expected := Eq{
			"str_col": Ptr("still a str"),
		}

		clause := ToClause(&filter2)
		assert.Equal(t, expected, clause.contents)
	})

	t.Run("Returns empty on nil input", func(t *testing.T) {
		expected := Eq{}

		clause := ToClause(nil)
		assert.Equal(t, expected, clause.contents)
	})
	t.Run("Returns empty on empty output", func(t *testing.T) {
		expected := Eq{}

		clause := ToClause(&thingyGetFilter{})
		assert.Equal(t, expected, clause.contents)
	})

	t.Run("Has an error if given a struct with no db tags", func(t *testing.T) {
		clause := ToClause(&thingyGetFilterWithNoTags{})
		assert.Error(t, clause.err)
	})

	t.Run("Has an error if given a struct with duplicate db tags", func(t *testing.T) {
		clause := ToClause(&thingyGetFilterWithDuplicateTags{
			Field1: Ptr("value1"),
			Field2: Ptr("value2"),
		})
		assert.Error(t, clause.err)
		assert.Equal(t, ErrDuplicateDBTags, clause.err)
	})
}

func TestToClauseAlias(t *testing.T) {
	f1 := thingyGetFilter{
		StrCol: Ptr("i am str"),
		IntCol: &[]int{100},
	}
	expected := Eq{
		"table.str_col": Ptr("i am str"),
		"table.int_col": &[]int{100},
	}

	clause := ToClauseAlias("table", &f1)
	assert.Equal(t, expected, clause.contents)
}
