package sqx

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type DeleteBuilder struct {
	builder sq.DeleteBuilder
	runner  sq.BaseRunner
	ctx     context.Context
	err     error
}

// ============================================
// BEGIN: squirrel-DeleteBuilder parity section
// ============================================

// Prefix adds an expression to the beginning of the query
func (b DeleteBuilder) Prefix(sql string, args ...interface{}) DeleteBuilder {
	return b.withBuilder(b.builder.Prefix(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b DeleteBuilder) PrefixExpr(expr sq.Sqlizer) DeleteBuilder {
	return b.withBuilder(b.builder.PrefixExpr(expr))
}

// From sets the table to be deleted from.
func (b DeleteBuilder) From(from string) DeleteBuilder {
	return b.withBuilder(b.builder.From(from))
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b DeleteBuilder) Where(pred interface{}, rest ...interface{}) DeleteBuilder {
	return b.withBuilder(b.builder.Where(pred, rest...))
}

// OrderBy adds ORDER BY expressions to the query.
func (b DeleteBuilder) OrderBy(orderBys ...string) DeleteBuilder {
	return b.withBuilder(b.builder.OrderBy(orderBys...))
}

// Limit sets a LIMIT clause on the query.
func (b DeleteBuilder) Limit(limit uint64) DeleteBuilder {
	return b.withBuilder(b.builder.Limit(limit))
}

// Offset sets a OFFSET clause on the query.
func (b DeleteBuilder) Offset(offset uint64) DeleteBuilder {
	return b.withBuilder(b.builder.Offset(offset))
}

// Suffix adds an expression to the end of the query
func (b DeleteBuilder) Suffix(sql string, args ...interface{}) DeleteBuilder {
	return b.withBuilder(b.builder.Suffix(sql, args...))
}

// ==========================================
// END: squirrel-UpdateBuilder parity section
// ==========================================

// Do executes the DeleteBuilder
func (b DeleteBuilder) Do() error {
	if b.err != nil {
		return b.err
	}
	_, err := b.builder.RunWith(b.runner).ExecContext(b.ctx)
	return err
}

// Debug prints the DeleteBuilder state out to the provided logger
func (b DeleteBuilder) Debug() DeleteBuilder {
	debug(b.ctx, b.builder)
	return b
}

func (b DeleteBuilder) withBuilder(builder sq.DeleteBuilder) DeleteBuilder {
	return DeleteBuilder{builder: builder, runner: b.runner, ctx: b.ctx, err: b.err}
}
