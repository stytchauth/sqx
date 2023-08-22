package sqx

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
)

type UpdateBuilder struct {
	builder    sq.UpdateBuilder
	runner     BaseRunner
	ctx        context.Context
	err        error
	hasChanges bool
}

// ============================================
// BEGIN: squirrel-UpdateBuilder parity section
// ============================================

// Prefix adds an expression to the beginning of the query
func (b UpdateBuilder) Prefix(sql string, args ...interface{}) UpdateBuilder {
	return b.withBuilder(b.builder.Prefix(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b UpdateBuilder) PrefixExpr(expr Sqlizer) UpdateBuilder {
	return b.withBuilder(b.builder.PrefixExpr(expr))
}

// Set adds SET clauses to the query.
func (b UpdateBuilder) Set(column string, value any) UpdateBuilder {
	return b.
		withBuilder(b.builder.Set(column, value)).
		withChanges()
}

// SetMap is a convenience method which calls Set for each key/value pair in clauses.
func (b UpdateBuilder) SetMap(clauses map[string]any, errors ...error) UpdateBuilder {
	for _, err := range errors {
		if err != nil {
			return b.withError(err)
		}
	}
	if len(clauses) == 0 {
		return b
	}
	return b.
		withBuilder(b.builder.SetMap(clauses)).
		withChanges()
}

// Where adds WHERE expressions to the query.
//
// See SelectBuilder.Where for more information.
func (b UpdateBuilder) Where(pred interface{}, rest ...interface{}) UpdateBuilder {
	return b.withBuilder(b.builder.Where(pred, rest...))
}

// OrderBy adds ORDER BY expressions to the query.
func (b UpdateBuilder) OrderBy(orderBys ...string) UpdateBuilder {
	return b.withBuilder(b.builder.OrderBy(orderBys...))
}

// Limit sets a LIMIT clause on the query.
func (b UpdateBuilder) Limit(limit uint64) UpdateBuilder {
	return b.withBuilder(b.builder.Limit(limit))
}

// Offset sets a OFFSET clause on the query.
func (b UpdateBuilder) Offset(offset uint64) UpdateBuilder {
	return b.withBuilder(b.builder.Offset(offset))
}

// Suffix adds an expression to the end of the query
func (b UpdateBuilder) Suffix(sql string, args ...interface{}) UpdateBuilder {
	return b.withBuilder(b.builder.Suffix(sql, args...))
}

// SuffixExpr adds an expression to the end of the query
func (b UpdateBuilder) SuffixExpr(expr Sqlizer) UpdateBuilder {
	return b.withBuilder(b.builder.SuffixExpr(expr))
}

// ==========================================
// END: squirrel-UpdateBuilder parity section
// ==========================================

// Do executes the UpdateBuilder
func (b UpdateBuilder) Do() error {
	if b.err != nil {
		return b.err
	}
	if !b.hasChanges {
		log.Println("Skipping write to DB - no updates set")
		return nil
	}
	_, err := b.builder.RunWith(b.runner).ExecContext(b.ctx)
	return err
}

// Debug prints the UpdateBuilder state out to the provided logger
func (b UpdateBuilder) Debug() UpdateBuilder {
	debug(b.ctx, b.builder)
	return b
}

func (b UpdateBuilder) withError(err error) UpdateBuilder {
	if b.err != nil {
		return b
	}
	return UpdateBuilder{builder: b.builder, runner: b.runner, ctx: b.ctx, err: err, hasChanges: b.hasChanges}
}

func (b UpdateBuilder) withBuilder(builder sq.UpdateBuilder) UpdateBuilder {
	return UpdateBuilder{builder: builder, runner: b.runner, ctx: b.ctx, err: b.err, hasChanges: b.hasChanges}
}

func (b UpdateBuilder) withChanges() UpdateBuilder {
	return UpdateBuilder{builder: b.builder, runner: b.runner, ctx: b.ctx, err: b.err, hasChanges: true}
}
