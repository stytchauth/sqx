package sqx_test

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stytchauth/sqx"
)

type Widget struct {
	ID      string  `db:"widget_id"`
	Status  string  `db:"status"`
	Enabled bool    `db:"enabled"`
	OwnerID *string `db:"owner_id"`
}

func (w Widget) toSetMap() (map[string]any, error) {
	if w.ID == "" {
		return nil, fmt.Errorf("missing ID")
	}
	if w.Status == "" {
		return nil, fmt.Errorf("missing Status")
	}

	return sqx.ToSetMap(&w)
}

type dbWidget struct{}

func newDBWidget() dbWidget {
	return dbWidget{}
}

func (d *dbWidget) Create(ctx context.Context, tx sqx.Queryable, w *Widget) error {
	return sqx.Write(ctx).
		WithQueryable(tx).
		Insert("sqx_widgets_test").
		SetMap(w.toSetMap()).
		Do()
}

func (d *dbWidget) CreateMany(ctx context.Context, tx sqx.Queryable, ws []Widget) error {
	return sqx.TypedWrite[Widget](ctx).
		WithQueryable(tx).
		InsertMany("sqx_widgets_test").
		FromItems(ws).
		Do()
}

func (d *dbWidget) Delete(ctx context.Context, tx sqx.Queryable, widgetID string) error {
	return sqx.Write(ctx).
		WithQueryable(tx).
		Delete("sqx_widgets_test").
		Where(sqx.Eq{"widget_id": widgetID}).
		Do()
}

type widgetUpdateFilter struct {
	Status  *string              `db:"status"`
	Enabled *bool                `db:"enabled"`
	OwnerID sqx.Nullable[string] `db:"owner_id"`
}

func (w *widgetUpdateFilter) toSetMap() (map[string]any, error) {
	if w.Status != nil && *w.Status == "Greasy" {
		return nil, fmt.Errorf("widgets cannot be greasy")
	}
	return sqx.ToSetMap(w)
}

func (d *dbWidget) Update(ctx context.Context, tx sqx.Queryable, widgetID string, f *widgetUpdateFilter) error {
	_, err := d.UpdateResult(ctx, tx, widgetID, f)
	return err
}

func (d *dbWidget) UpdateResult(ctx context.Context, tx sqx.Queryable, widgetID string, f *widgetUpdateFilter) (sql.Result, error) {
	return sqx.Write(ctx).
		WithQueryable(tx).
		Update("sqx_widgets_test").
		Where(sqx.Eq{"widget_id": widgetID}).
		SetMap(f.toSetMap()).
		DoResult()
}

func (d *dbWidget) GetByID(ctx context.Context, tx sqx.Queryable, widgetID string) (*Widget, error) {
	return sqx.Read[Widget](ctx).
		WithQueryable(tx).
		Select("*").
		From("sqx_widgets_test").
		Where(sqx.Eq{"widget_id": widgetID}).
		OneStrict()
}

type widgetGetFilter struct {
	WidgetID *[]string `db:"widget_id"`
	Status   *string   `db:"status"`
}

func (d *dbWidget) Get(ctx context.Context, tx sqx.Queryable, f *widgetGetFilter) ([]Widget, error) {
	return sqx.Read[Widget](ctx).
		WithQueryable(tx).
		Select("*").
		From("sqx_widgets_test").
		Where(sqx.ToClause(f)).
		All()
}

func (d *dbWidget) First(ctx context.Context, tx sqx.Queryable, f *widgetGetFilter) (*Widget, error) {
	return sqx.Read[Widget](ctx).
		WithQueryable(tx).
		Select("*").
		From("sqx_widgets_test").
		Where(sqx.ToClause(f)).
		First()
}

func (d *dbWidget) GetAll(ctx context.Context, tx sqx.Queryable) ([]Widget, error) {
	return sqx.Read[Widget](ctx).
		WithQueryable(tx).
		Select("*").
		From("sqx_widgets_test").
		OrderBy("widget_id DESC").
		All()
}
