package sqx

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	scan "github.com/blockloop/scan/v2"
)

// SelectBuilder wraps squirrel.SelectBuilder and adds syntactic sugar for
// common usage patterns
type SelectBuilder[T any] struct {
	builder sq.SelectBuilder
	runner  sq.BaseRunner
	ctx     context.Context
	err     error
}

func (s SelectBuilder[T]) One() (*T, error) {
	dest, err := s.All()

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

func (s SelectBuilder[T]) OneScalar() (T, error) {
	ptr, err := s.One()
	if err != nil {
		var uninitialized T
		return uninitialized, err
	}
	return *ptr, nil
}

func (s SelectBuilder[T]) All() ([]T, error) {
	rows, err := s.query()
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

func (s SelectBuilder[T]) query() (*sql.Rows, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.runner == nil {
		return nil, errors.New("no runner")
	}
	if s.ctx == nil {
		return nil, errors.New("no ctx")
	}
	return s.builder.RunWith(s.runner).QueryContext(s.ctx)
}

// TODO: Whenever you want to use a method on SelectBuilder[T] that isn't defined, implement it here
// TODO: Can we code gen this?
// Example - Having
//func (s SelectBuilder[T]) Having(pred interface{}, rest ...interface{}) SelectBuilder[T] {
//  return s.withBuilder(s.builder.Having(pred, rest...))
//}

func (s SelectBuilder[T]) From(from string) SelectBuilder[T] {
	return s.withBuilder(s.builder.From(from))
}

func (s SelectBuilder[T]) FromSelect(from SelectBuilder[T], alias string) SelectBuilder[T] {
	return s.withBuilder(s.builder.FromSelect(from.builder, alias))
}

func (s SelectBuilder[T]) Join(join string, rest ...interface{}) SelectBuilder[T] {
	return s.withBuilder(s.builder.Join(join, rest...))
}

func (s SelectBuilder[T]) LeftJoin(join string, rest ...interface{}) SelectBuilder[T] {
	return s.withBuilder(s.builder.LeftJoin(join, rest...))
}

func (s SelectBuilder[T]) Where(pred interface{}, rest ...interface{}) SelectBuilder[T] {
	return s.withBuilder(s.builder.Where(pred, rest...))
}

func (s SelectBuilder[T]) OrderBy(orderBys ...string) SelectBuilder[T] {
	return s.withBuilder(s.builder.OrderBy(orderBys...))
}

func (s SelectBuilder[T]) Columns(columns ...string) SelectBuilder[T] {
	return s.withBuilder(s.builder.Columns(columns...))
}

func (s SelectBuilder[T]) Limit(limit uint64) SelectBuilder[T] {
	return s.withBuilder(s.builder.Limit(limit))
}

func (s SelectBuilder[T]) GroupBy(groupBys ...string) SelectBuilder[T] {
	return s.withBuilder(s.builder.GroupBy(groupBys...))
}

func (s SelectBuilder[T]) Suffix(sql string, rest ...interface{}) SelectBuilder[T] {
	return s.withBuilder(s.builder.Suffix(sql, rest...))
}

func (s SelectBuilder[T]) UnionAll(other SelectBuilder[T]) SelectBuilder[T] {
	query, args, err := other.builder.ToSql()
	if err != nil {
		return s.withError(err)
	}
	return s.withBuilder(s.builder.Suffix("UNION ALL ("+query+")", args...))
}

func (s SelectBuilder[T]) Using(operation func(builder SelectBuilder[T]) SelectBuilder[T]) SelectBuilder[T] {
	return operation(s)
}

func (s SelectBuilder[T]) Debug() SelectBuilder[T] {
	debug(s.ctx, s.builder)
	return s
}

func (s SelectBuilder[T]) withBuilder(builder sq.SelectBuilder) SelectBuilder[T] {
	return SelectBuilder[T]{builder: builder, runner: s.runner, ctx: s.ctx, err: s.err}
}

func (s SelectBuilder[T]) withError(err error) SelectBuilder[T] {
	return SelectBuilder[T]{builder: s.builder, runner: s.runner, ctx: s.ctx, err: err}
}
