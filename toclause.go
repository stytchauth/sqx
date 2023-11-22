package sqx

import (
	"errors"

	scan "github.com/blockloop/scan/v2"
)

var NoDBTagsError = errors.New("No db tags detected")

// Clause stores an Eq result, but also holds an error if one occurred during the conversion. You may think of this
// struct as a (Eq, error) tuple that implements the Sqlizer interface.
type Clause struct {
	contents Eq
	err      error
}

// ToSql calls the underlying Eq's ToSql method, but returns the error if one occurred when the Clause was constructed.
func (c *Clause) ToSql() (string, []interface{}, error) {
	if c.err != nil {
		return "", nil, c.err
	}
	return c.contents.ToSql()
}

// ToClause converts a filter interface to a SQL Where clause by introspecting its db tags
func ToClause(v any, excluded ...string) *Clause {
	if isNil(v) {
		return &Clause{contents: Eq{}, err: nil}
	}
	cols, err := scan.ColumnsStrict(v, excluded...)
	if err != nil {
		return &Clause{contents: nil, err: err}
	}
	if len(cols) == 0 {
		return &Clause{contents: nil, err: NoDBTagsError}
	}
	vals, err := scan.Values(cols, v)
	if err != nil {
		return &Clause{contents: nil, err: err}
	}
	contents := Eq{}
	for i := range cols {
		if !isNil(vals[i]) {
			contents[cols[i]] = vals[i]
		}
	}
	return &Clause{contents: contents, err: nil}
}

// ToClauseAlias is like ToClause, but takes in a table alias
func ToClauseAlias(tableName string, v any, excluded ...string) *Clause {
	clause := ToClause(v, excluded...)
	if clause.err != nil {
		return clause
	}
	aliasedClause := &Clause{contents: Eq{}, err: nil}
	for key, value := range clause.contents {
		aliasedClause.contents[tableName+"."+key] = value
	}
	return aliasedClause
}
