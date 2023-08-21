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

func (s DeleteBuilder) Where(pred interface{}, rest ...interface{}) DeleteBuilder {
	return s.withBuilder(s.builder.Where(pred, rest...))
}

func (s DeleteBuilder) Do() error {
	if s.err != nil {
		return s.err
	}
	_, err := s.builder.RunWith(s.runner).ExecContext(s.ctx)
	return err
}

func (s DeleteBuilder) Debug() DeleteBuilder {
	debug(s.ctx, s.builder)
	return s
}

func (s DeleteBuilder) withBuilder(builder sq.DeleteBuilder) DeleteBuilder {
	return DeleteBuilder{builder: builder, runner: s.runner, ctx: s.ctx, err: s.err}
}
