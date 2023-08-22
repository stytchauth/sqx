package sqx_test

import (
	"context"
	"fmt"

	"github.com/stytchauth/sqx"
)

type Widget struct {
	ID      string `db:"widget_id"`
	Status  string `db:"Status"`
	Enabled bool   `db:"enabled"`
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

type dbWidget struct {
	db sqx.Queryable
}

func newDBWidget(db sqx.Queryable) dbWidget {
	return dbWidget{db: db}
}

func (d *dbWidget) Create(ctx context.Context, w *Widget) error {
	return sqx.Write(ctx).
		Insert("sqx_widgets_test").
		SetMap(w.toSetMap()).
		Debug().
		Do()
}

type widgetUpdateFilter struct {
	Status  *string `db:"status"`
	Enabled *bool   `db:"enabled"`
}

func (w *widgetUpdateFilter) toSetMap() (map[string]any, error) {
	if w.Status != nil && *w.Status == "Greasy" {
		return nil, fmt.Errorf("widgets cannot be greasy")
	}
	return sqx.ToSetMap(w)
}

func (d *dbWidget) Update(ctx context.Context, widgetID string, f *widgetUpdateFilter) error {
	if !sqx.ContainsUpdates(f) {
		return nil
	}
	return sqx.Write(ctx).
		Update("sqx_widgets_test").
		Where(sqx.Eq{"widget_id": widgetID}).
		SetMap(f.toSetMap()).
		Debug().
		Do()
}

func (d *dbWidget) GetByID(ctx context.Context, widgetID string) (*Widget, error) {
	return sqx.Read[Widget](ctx).
		Select("*").
		From("sqx_widgets_test").
		Where(sqx.Eq{"widget_id": widgetID}).
		Debug().
		One()
}

type widgetGetFilter struct {
	widgetID *[]string `db:"id"`
}

func (d *dbWidget) Get(ctx context.Context, f *widgetGetFilter) ([]Widget, error) {
	return sqx.Read[Widget](ctx).
		Select("*").
		From("sqx_widgets_test").
		Where(sqx.ToClause(f)).
		Debug().
		All()
}

func (d *dbWidget) GetAll(ctx context.Context) ([]Widget, error) {
	return sqx.Read[Widget](ctx).
		Select("*").
		Debug().
		From("sqx_widgets_test").
		Debug().
		OrderBy("widget_id DESC").
		Debug().
		All()
}
