package sqx

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type InsertBuilder struct {
	builder   sq.InsertBuilder
	queryable Queryable
	ctx       context.Context
	err       error
	logger    Logger
}

// ============================================
// BEGIN: squirrel-InsertBuilder parity section
// ============================================

// Prefix adds an expression to the beginning of the query
func (b InsertBuilder) Prefix(sql string, args ...interface{}) InsertBuilder {
	return b.withBuilder(b.builder.Prefix(sql, args...))
}

// PrefixExpr adds an expression to the very beginning of the query
func (b InsertBuilder) PrefixExpr(expr Sqlizer) InsertBuilder {
	return b.withBuilder(b.builder.PrefixExpr(expr))
}

// Options adds keyword options before the INTO clause of the query.
func (b InsertBuilder) Options(options ...string) InsertBuilder {
	return b.withBuilder(b.builder.Options(options...))
}

// Columns adds insert columns to the query.
func (b InsertBuilder) Columns(columns ...string) InsertBuilder {
	return b.withBuilder(b.builder.Columns(columns...))
}

// Values adds a single row's values to the query.
func (b InsertBuilder) Values(values ...any) InsertBuilder {
	return b.withBuilder(b.builder.Values(values...))
}

// Suffix adds an expression to the end of the query
func (b InsertBuilder) Suffix(sql string, args ...interface{}) InsertBuilder {
	return b.withBuilder(b.builder.Suffix(sql, args...))
}

// SuffixExpr adds an expression to the end of the query
func (b InsertBuilder) SuffixExpr(expr Sqlizer) InsertBuilder {
	return b.withBuilder(b.builder.SuffixExpr(expr))
}

// SetMap set columns and values for insert builder from a map of column name and value
// note that it will reset all previous columns and values was set if any
func (b InsertBuilder) SetMap(clauses map[string]interface{}, errors ...error) InsertBuilder {
	for _, err := range errors {
		if err != nil {
			return b.withError(err)
		}
	}
	return b.withBuilder(b.builder.SetMap(clauses))
}

// ==========================================
// END: squirrel-InsertBuilder parity section
// ==========================================

// Do executes the InsertBuilder
func (b InsertBuilder) Do() error {
	if b.err != nil {
		return b.err
	}
	if b.queryable == nil {
		return fmt.Errorf("missing queryable - call SetDefaultQueryable or WithQueryable to set it")
	}
	_, err := b.builder.RunWith(runShim{b.queryable}).ExecContext(b.ctx)
	return err
}

// Debug prints the InsertBuilder state out to the provided logger
func (b InsertBuilder) Debug() InsertBuilder {
	debug(b.logger, b.builder)
	return b
}

// WithQueryable configures a Queryable for this InsertBuilder instance
func (b InsertBuilder) WithQueryable(queryable Queryable) InsertBuilder {
	return InsertBuilder{builder: b.builder, queryable: queryable, logger: b.logger, ctx: b.ctx, err: b.err}
}

// WithLogger configures a Queryable for this InsertBuilder instance
func (b InsertBuilder) WithLogger(logger Logger) InsertBuilder {
	return InsertBuilder{builder: b.builder, queryable: b.queryable, logger: logger, ctx: b.ctx, err: b.err}
}

func (b InsertBuilder) withError(err error) InsertBuilder {
	if b.err != nil {
		return b
	}
	return InsertBuilder{builder: b.builder, queryable: b.queryable, logger: b.logger, ctx: b.ctx, err: err}
}

func (b InsertBuilder) withBuilder(builder sq.InsertBuilder) InsertBuilder {
	return InsertBuilder{builder: builder, queryable: b.queryable, logger: b.logger, ctx: b.ctx, err: b.err}
}
