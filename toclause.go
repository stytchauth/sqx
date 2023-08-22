package sqx

import (
	sq "github.com/Masterminds/squirrel"
	scan "github.com/blockloop/scan"
)

type Clause struct {
	contents sq.Eq
	err      error
}

//nolint:stylecheck
func (c *Clause) ToSql() (string, []interface{}, error) {
	if c.err != nil {
		return "", nil, c.err
	}
	return c.contents.ToSql()
}

// ToClause converts a filter interface to a SQL Where clause by introspecting its db tags
func ToClause(v any, excluded ...string) *Clause {
	if isNil(v) {
		return &Clause{contents: sq.Eq{}, err: nil}
	}
	cols, err := scan.ColumnsStrict(v, excluded...)
	if err != nil {
		return &Clause{contents: nil, err: err}
	}
	vals, err := scan.Values(cols, v)
	if err != nil {
		return &Clause{contents: nil, err: err}
	}
	contents := sq.Eq{}
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
	aliasedClause := &Clause{contents: sq.Eq{}, err: nil}
	for key, value := range clause.contents {
		aliasedClause.contents[tableName+"."+key] = value
	}
	return aliasedClause
}
