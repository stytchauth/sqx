package sqx

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/blockloop/scan"
)

// SelectBuilder wraps squirrel.SelectBuilder and adds syntactic sugar for
// common usage patterns
type SelectBuilder[T any] struct {
	builder   sq.SelectBuilder
	queryable Queryable
	ctx       context.Context
	err       error
	logger    Logger
}

// ============================================
// BEGIN: squirrel-SelectBuilder parity section
// ============================================

// Prefix adds an expression to the beginning of the query
func (b SelectBuilder[T]) Prefix(sql string, args ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.Prefix(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b SelectBuilder[T]) PrefixExpr(expr Sqlizer) SelectBuilder[T] {
	return b.withBuilder(b.builder.PrefixExpr(expr))
}

// Distinct adds a DISTINCT clause to the query.
func (b SelectBuilder[T]) Distinct() SelectBuilder[T] {
	return b.withBuilder(b.builder.Distinct())
}

// Options adds select option to the query
func (b SelectBuilder[T]) Options(options ...string) SelectBuilder[T] {
	return b.withBuilder(b.builder.Options(options...))
}

// Columns adds result columns to the query.
func (b SelectBuilder[T]) Columns(columns ...string) SelectBuilder[T] {
	return b.withBuilder(b.builder.Columns(columns...))
}

// RemoveColumns remove all columns from query.
// Must add a new column with Column or Columns methods, otherwise
// return a error.
func (b SelectBuilder[T]) RemoveColumns() SelectBuilder[T] {
	return b.withBuilder(b.builder.RemoveColumns())
}

// Column adds a result column to the query.
// Unlike Columns, Column accepts args which will be bound to placeholders in
// the columns string, for example:
//
//	Column("IF(col IN ("+squirrel.Placeholders(3)+"), 1, 0) as col", 1, 2, 3)
func (b SelectBuilder[T]) Column(column any, args ...any) SelectBuilder[T] {
	return b.withBuilder(b.builder.Column(column, args...))
}

// From sets the FROM clause of the query.
func (b SelectBuilder[T]) From(from string) SelectBuilder[T] {
	return b.withBuilder(b.builder.From(from))
}

// FromSelect sets a subquery into the FROM clause of the query.
func (b SelectBuilder[T]) FromSelect(from SelectBuilder[T], alias string) SelectBuilder[T] {
	return b.withBuilder(b.builder.FromSelect(from.builder, alias))
}

// JoinClause adds a join clause to the query.
func (b SelectBuilder[T]) JoinClause(pred interface{}, args ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.JoinClause(pred, args...))
}

// Join adds a JOIN clause to the query.
func (b SelectBuilder[T]) Join(join string, rest ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.Join(join, rest...))
}

// LeftJoin adds a LEFT JOIN clause to the query.
func (b SelectBuilder[T]) LeftJoin(join string, rest ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.LeftJoin(join, rest...))
}

// RightJoin adds a RIGHT JOIN clause to the query.
func (b SelectBuilder[T]) RightJoin(join string, rest ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.RightJoin(join, rest...))
}

// InnerJoin adds a INNER JOIN clause to the query.
func (b SelectBuilder[T]) InnerJoin(join string, rest ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.InnerJoin(join, rest...))
}

// CrossJoin adds a CROSS JOIN clause to the query.
func (b SelectBuilder[T]) CrossJoin(join string, rest ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.CrossJoin(join, rest...))
}

// Where adds an expression to the WHERE clause of the query.
//
// Expressions are ANDed together in the generated SQL.
//
// Where accepts several types for its pred argument:
//
// nil OR "" - ignored.
//
// string - SQL expression.
// If the expression has SQL placeholders then a set of arguments must be passed
// as well, one for each placeholder.
//
// map[string]interface{} OR Eq - map of SQL expressions to values. Each key is
// transformed into an expression like "<key> = ?", with the corresponding value
// bound to the placeholder. If the value is nil, the expression will be "<key>
// IS NULL". If the value is an array or slice, the expression will be "<key> IN
// (?,?,...)", with one placeholder for each item in the value. These expressions
// are ANDed together.
//
// Where will panic if pred isn't any of the above types.
func (b SelectBuilder[T]) Where(pred interface{}, rest ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.Where(pred, rest...))
}

// GroupBy adds GROUP BY expressions to the query.
func (b SelectBuilder[T]) GroupBy(groupBys ...string) SelectBuilder[T] {
	return b.withBuilder(b.builder.GroupBy(groupBys...))
}

// Having adds an expression to the HAVING clause of the query.
//
// See Where.
func (b SelectBuilder[T]) Having(pred interface{}, rest ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.Having(pred, rest...))
}

// OrderByClause adds ORDER BY clause to the query.
func (b SelectBuilder[T]) OrderByClause(pred interface{}, args ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.OrderByClause(pred, args...))
}

// OrderBy adds ORDER BY expressions to the query.
func (b SelectBuilder[T]) OrderBy(orderBys ...string) SelectBuilder[T] {
	return b.withBuilder(b.builder.OrderBy(orderBys...))
}

// Limit sets a LIMIT clause on the query.
func (b SelectBuilder[T]) Limit(limit uint64) SelectBuilder[T] {
	return b.withBuilder(b.builder.Limit(limit))
}

// RemoveLimit removes LIMIT clause
func (b SelectBuilder[T]) RemoveLimit() SelectBuilder[T] {
	return b.withBuilder(b.builder.RemoveLimit())
}

// Offset sets a OFFSET clause on the query.
func (b SelectBuilder[T]) Offset(offset uint64) SelectBuilder[T] {
	return b.withBuilder(b.builder.Offset(offset))
}

// RemoveOffset removes OFFSET clause.
func (b SelectBuilder[T]) RemoveOffset() SelectBuilder[T] {
	return b.withBuilder(b.builder.RemoveOffset())
}

// Suffix adds an expression to the end of the query
func (b SelectBuilder[T]) Suffix(sql string, rest ...interface{}) SelectBuilder[T] {
	return b.withBuilder(b.builder.Suffix(sql, rest...))
}

// SuffixExpr adds an expression to the end of the query
func (b SelectBuilder[T]) SuffixExpr(expr Sqlizer) SelectBuilder[T] {
	return b.withBuilder(b.builder.SuffixExpr(expr))
}

// ==========================================
// END: squirrel-SelectBuilder parity section
// ==========================================

func (b SelectBuilder[T]) UnionAll(other SelectBuilder[T]) SelectBuilder[T] {
	query, args, err := other.builder.ToSql()
	if err != nil {
		return b.withError(err)
	}
	return b.withBuilder(b.builder.Suffix("UNION ALL ("+query+")", args...))
}

func (b SelectBuilder[T]) One() (*T, error) {
	dest, err := b.All()

	if err != nil {
		return nil, err
	} else if len(dest) == 0 {
		// Since we called RowsStrict (plural) above, no `sql.ErrNoRows` would have been raised
		// since a slice of zero elements is a valid return value. So we raise it ourselves now.
		return nil, sql.ErrNoRows
	}
	// TODO: Optionally warn if len(dest) > 1
	return &dest[0], nil
}

func (b SelectBuilder[T]) OneScalar() (T, error) {
	ptr, err := b.One()
	if err != nil {
		var uninitialized T
		return uninitialized, err
	}
	return *ptr, nil
}

func (b SelectBuilder[T]) All() ([]T, error) {
	rows, err := b.query()
	if err != nil {
		return nil, err
	}

	var dest []T
	err = scan.RowsStrict(&dest, rows)

	if err != nil {
		return nil, err
	} else {
		return dest, nil
	}
}

func (b SelectBuilder[T]) query() (*sql.Rows, error) {
	if b.err != nil {
		return nil, b.err
	}
	if b.queryable == nil {
		return nil, errors.New("no queryable")
	}
	if b.ctx == nil {
		return nil, errors.New("no ctx")
	}
	if b.queryable == nil {
		return nil, fmt.Errorf("missing queryable - call SetDefaultQueryable or WithQueryable to set it")
	}
	return b.builder.RunWith(runShim{b.queryable}).QueryContext(b.ctx)
}
func (b SelectBuilder[T]) Debug() SelectBuilder[T] {
	debug(b.logger, b.builder)
	return b
}

// WithQueryable configures a Queryable for this SelectBuilder instance
func (b SelectBuilder[T]) WithQueryable(queryable Queryable) SelectBuilder[T] {
	return SelectBuilder[T]{builder: b.builder, queryable: queryable, logger: b.logger, ctx: b.ctx, err: b.err}
}

// WithLogger configures a Queryable for this SelectBuilder instance
func (b SelectBuilder[T]) WithLogger(logger Logger) SelectBuilder[T] {
	return SelectBuilder[T]{builder: b.builder, queryable: b.queryable, logger: logger, ctx: b.ctx, err: b.err}
}

func (b SelectBuilder[T]) withBuilder(builder sq.SelectBuilder) SelectBuilder[T] {
	return SelectBuilder[T]{builder: builder, queryable: b.queryable, logger: b.logger, ctx: b.ctx, err: b.err}
}

func (b SelectBuilder[T]) withError(err error) SelectBuilder[T] {
	return SelectBuilder[T]{builder: b.builder, queryable: b.queryable, logger: b.logger, ctx: b.ctx, err: err}
}
