package sqx

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
)

type UpdateBuilder struct {
	builder    sq.UpdateBuilder
	runner     sq.BaseRunner
	ctx        context.Context
	err        error
	hasChanges bool
}

func (s UpdateBuilder) Set(column string, value any) UpdateBuilder {
	return s.
		withBuilder(s.builder.Set(column, value)).
		withChanges()
}

func (s UpdateBuilder) SetMap(clauses map[string]interface{}, errors ...error) UpdateBuilder {
	for _, err := range errors {
		if err != nil {
			return s.withError(err)
		}
	}
	if len(clauses) == 0 {
		return s
	}
	return s.
		withBuilder(s.builder.SetMap(clauses)).
		withChanges()
}

func (s UpdateBuilder) Where(pred interface{}, rest ...interface{}) UpdateBuilder {
	return s.withBuilder(s.builder.Where(pred, rest...))
}

func (s UpdateBuilder) Using(operation func(builder UpdateBuilder) UpdateBuilder) UpdateBuilder {
	return operation(s)
}

func (s UpdateBuilder) Debug() UpdateBuilder {
	debug(s.ctx, s.builder)
	return s
}

func (s UpdateBuilder) Do() error {
	if s.err != nil {
		return s.err
	}
	if !s.hasChanges {
		log.Println("Skipping write to DB - no updates set")
		return nil
	}
	_, err := s.builder.RunWith(s.runner).ExecContext(s.ctx)
	return err
}

func (s UpdateBuilder) withError(err error) UpdateBuilder {
	if s.err != nil {
		return s
	}
	return UpdateBuilder{builder: s.builder, runner: s.runner, ctx: s.ctx, err: err, hasChanges: s.hasChanges}
}

func (s UpdateBuilder) withBuilder(builder sq.UpdateBuilder) UpdateBuilder {
	return UpdateBuilder{builder: builder, runner: s.runner, ctx: s.ctx, err: s.err, hasChanges: s.hasChanges}
}

func (s UpdateBuilder) withChanges() UpdateBuilder {
	return UpdateBuilder{builder: s.builder, runner: s.runner, ctx: s.ctx, err: s.err, hasChanges: true}
}
