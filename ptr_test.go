package sqx

import "fmt"

func ExamplePtr() {
	type dogFilter struct {
		Breed      *string `db:"breed"`
		PlaysFetch *bool   `db:"plays_fetch"`
	}
	clause := ToClause(&dogFilter{
		Breed: Ptr("husky"),
	})
	sql, args, _ := clause.ToSql()
	fmt.Printf("%s, %s", sql, args)

	// Output:
	// breed = ?, [husky]
}
