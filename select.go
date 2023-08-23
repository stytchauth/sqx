package sqx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/blockloop/scan/v2"
	sq "github.com/stytchauth/squirrel"
)

// SelectBuilder wraps squirrel.SelectBuilder and adds syntactic sugar for common usage patterns.
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

// UnionAll adds a UNION ALL clause to the query from another SelectBuilder of the same type.
func (b SelectBuilder[T]) UnionAll(other SelectBuilder[T]) SelectBuilder[T] {
	query, args, err := other.builder.ToSql()
	if err != nil {
		return b.withError(err)
	}
	return b.withBuilder(b.builder.Suffix("UNION ALL ("+query+")", args...))
}

// one returns a single result from the query, or an error if there was a problem. It may be run in strict or non-strict
// mode. In non-strict mode, a warning is logged if more than one result is returned in the query. In strict mode, this
// turns into an ErrTooManyRows error. If the underlying query is *expected* to return more than one row and this is not
// cause for concern, you should instead use First.
func (b SelectBuilder[T]) one(strict bool) (*T, error) {
	dest, err := b.All()

	if err != nil {
		return nil, err
	} else if len(dest) == 0 {
		// Since we called RowsStrict (plural) above, no `sql.ErrNoRows` would have been raised
		// since a slice of zero elements is a valid return value. So we raise it ourselves now.
		return nil, sql.ErrNoRows
	}

	if len(dest) > 1 {
		if strict {
			return nil, ErrTooManyRows{Expected: 1, Actual: len(dest)}
		} else if b.logger != nil {
			b.logger.Printf("[WARN] sqx: in call to One, got %d rows, returning first result", len(dest))
		}
	}

	return &dest[0], nil
}

// One returns a single result from the query, or an error if there was a problem. This runs in "non-strict" mode which
// means that if the underlying query returns more than one row, a warning is logged but no error is raised. If you want
// to raise an error if the underlying query returns more than one result, use OneStrict. If you instead expect that
// more than one result may be returned and this is not cause for concern, use First.
func (b SelectBuilder[T]) One() (*T, error) {
	return b.one(false)
}

// OneStrict returns a single result from the query, or an error if there was a problem. This runs in "strict" mode
// which means that if the underlying query returns more than one row, an error is raised. You may instead use One to
// downgrade this error into a warning from the saved logger, or First to return the first result for cases where you
// expect more than one result can be returned from the underlying query and this is not cause for concern.
func (b SelectBuilder[T]) OneStrict() (*T, error) {
	return b.one(true)
}

// OneScalar is like One but dereferences the result into a scalar value. If an error is raised, the scalar value will
// be the zero value of the type.
func (b SelectBuilder[T]) OneScalar() (T, error) {
	ptr, err := b.One()
	if err != nil {
		var uninitialized T
		return uninitialized, err
	}
	return *ptr, nil
}

// OneScalarStrict is like OneStrict but dereferences the result into a scalar value. If an error is raised, the scalar
// value will be the zero value of the type.
func (b SelectBuilder[T]) OneScalarStrict() (T, error) {
	ptr, err := b.OneStrict()
	if err != nil {
		var uninitialized T
		return uninitialized, err
	}
	return *ptr, nil
}

// First returns the first result from the query, or an error if there was a problem. This is useful for queries that
// are expected to return more than one result, but you only care about the first one. Note that if you haven't added an
// ORDER BY clause to your query, the first result is not guaranteed to be the same each time you run the query.
func (b SelectBuilder[T]) First() (*T, error) {
	rows, err := b.query()
	if err != nil {
		return nil, err
	}

	var dest *T
	err = scan.RowStrict(dest, rows)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

// FirstScalar is like First but dereferences the result into a scalar value. If an error is raised, the scalar value
// will be the zero value of the type.
func (b SelectBuilder[T]) FirstScalar() (T, error) {
	ptr, err := b.First()
	if err != nil {
		var uninitialized T
		return uninitialized, err
	}
	return *ptr, nil
}

// All returns all results from the query as a slice of T.
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

// Debug prints the SQL query using the builder's logger and then returns b, unmodified. If the builder has no logger
// set (and SetDefaultLogger has not been called), then log.Printf is used instead.
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
