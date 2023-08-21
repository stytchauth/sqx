package sqx

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type InsertBuilder struct {
	builder sq.InsertBuilder
	runner  sq.BaseRunner
	ctx     context.Context
	err     error
}

func (s InsertBuilder) SetMap(clauses map[string]interface{}, errors ...error) InsertBuilder {
	for _, err := range errors {
		if err != nil {
			return s.withError(err)
		}
	}
	return s.withBuilder(s.builder.SetMap(clauses))
}

func (s InsertBuilder) Columns(columns ...string) InsertBuilder {
	return s.withBuilder(s.builder.Columns(columns...))
}

func (s InsertBuilder) Values(values ...any) InsertBuilder {
	return s.withBuilder(s.builder.Values(values...))
}

func (s InsertBuilder) Suffix(sql string, args ...any) InsertBuilder {
	return s.withBuilder(s.builder.Suffix(sql, args...))
}

func (s InsertBuilder) SuffixExpr(expr sq.Sqlizer) InsertBuilder {
	return s.withBuilder(s.builder.SuffixExpr(expr))
}

func (s InsertBuilder) Do() error {
	if s.err != nil {
		return s.err
	}
	_, err := s.builder.RunWith(s.runner).ExecContext(s.ctx)
	return err
}

func (s InsertBuilder) Debug() InsertBuilder {
	debug(s.ctx, s.builder)
	return s
}

func (s InsertBuilder) withError(err error) InsertBuilder {
	if s.err != nil {
		return s
	}
	return InsertBuilder{builder: s.builder, runner: s.runner, ctx: s.ctx, err: err}
}

func (s InsertBuilder) withBuilder(builder sq.InsertBuilder) InsertBuilder {
	return InsertBuilder{builder: builder, runner: s.runner, ctx: s.ctx, err: s.err}
}
