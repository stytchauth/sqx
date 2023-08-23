package sqx_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stytchauth/sqx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	sqx.SetDefaultLogger(sqx.MakeLogger(log.Printf))
}

func CreateDatabase() (*sql.DB, error) {
	return sql.Open("mysql", "sqx:sqx@tcp(localhost:4306)/sqx?parseTime=true")
}

// DB opens a new database connection for the duration of the test.
//
// If you don't need to persist data (e.g., for chained-request testing), use Tx instead.
func DB(t *testing.T) *sql.DB {
	db, err := CreateDatabase()
	if err != nil {
		t.Fatalf("Create database: %s", err.Error())
	}

	sqx.SetDefaultQueryable(db)
	t.Cleanup(func() {
		sqx.SetDefaultQueryable(nil)
		if err := db.Close(); err != nil {
			t.Logf("Close DB connection: %s", err.Error())
		}
	})
	return db
}

// Tx opens a new database transaction for use in tests. The transaction will be rolled back as
// part of test cleanup.
func Tx(t *testing.T) *sql.Tx {
	// From the test's perspective, this transaction _is_ the database, and that doesn't go away
	// within a time limit.
	ctx := context.Background()

	tx, err := DB(t).BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("DB connection failed: %s. Did you remember to run make services?", err.Error())
	}

	// This overwrites the default queryable for the entire package.
	// You probably don't want to do this - you probably want to use WithQueryable(tx) instead
	sqx.SetDefaultQueryable(tx)
	t.Cleanup(func() {
		sqx.SetDefaultQueryable(nil)
		if err := tx.Rollback(); err != nil {
			t.Logf("Could not roll back transaction: %s", err.Error())
			// Don't fail the test here, because it's (probably) not the test's fault. Although
			// the transaction might stay open "forever", we'll likely shut down the test DB
			// within a few minutes.
		}
	})
	return tx
}

func setupTestWidgetsTable(t *testing.T) *sql.Tx {
	db := Tx(t)
	_, err := db.Exec(`DROP TABLE IF EXISTS sqx_widgets_test;`)
	require.NoError(t, err)
	_, err = db.Exec(`
		CREATE TABLE sqx_widgets_test (
			widget_id		VARCHAR(128) NOT NULL,
			status			VARCHAR(128) NOT NULL,
			enabled			BOOLEAN NOT NULL,
			owner_id 		VARCHAR(128)
		)
	`)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := db.Exec(`DROP TABLE IF EXISTS sqx_widgets_test;`)
		require.NoError(t, err)
	})
	return db
}

func newWidget(status string) Widget {
	return Widget{
		ID:      uuid.New().String(),
		Status:  status,
		Enabled: true,
	}
}

func TestRead(t *testing.T) {
	setupTestWidgetsTable(t)
	dbWidget := newDBWidget()
	ctx := context.Background()

	w1 := newWidget("great")
	w2 := newWidget("fine")

	require.NoError(t, dbWidget.Create(ctx, &w1))
	require.NoError(t, dbWidget.Create(ctx, &w2))

	t.Run("Can read a single widget", func(t *testing.T) {
		w1db, err := dbWidget.GetByID(ctx, w1.ID)
		assert.NoError(t, err)
		assert.Equal(t, &w1, w1db)
	})

	t.Run("Can read multiple widgets", func(t *testing.T) {
		widgets, err := dbWidget.GetAll(ctx)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []Widget{w1, w2}, widgets)
	})

	t.Run("Can read multiple widgets using a filter", func(t *testing.T) {
		widgets, err := dbWidget.Get(ctx, &widgetGetFilter{
			WidgetID: &[]string{w1.ID, w2.ID},
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, []Widget{w1, w2}, widgets)
	})

	t.Run("Bubbles up errors from the DB exec", func(t *testing.T) {
		w1db, err := dbWidget.GetByID(ctx, "bad-id")
		assert.EqualError(t, sql.ErrNoRows, err.Error())
		assert.Nil(t, w1db)
	})

	t.Run("Raises error when too many rows returned in OneStrict", func(t *testing.T) {
		w3 := newWidget("alright")
		w3.ID = w1.ID
		require.NoError(t, dbWidget.Create(ctx, &w3))

		expected := sqx.ErrTooManyRows{Expected: 1, Actual: 2}
		_, err := dbWidget.GetByID(ctx, w1.ID)
		assert.EqualError(t, expected, err.Error())
	})
}

func TestInsert(t *testing.T) {
	ctx := context.Background()

	w1 := newWidget("great")

	// The happy path cases are already tested in TestRead
	// here are only failure cases!

	t.Run("Returns an error when SetMap fails", func(t *testing.T) {
		setupTestWidgetsTable(t)
		dbWidget := newDBWidget()
		// Creating an empty widget should not work
		err := dbWidget.Create(ctx, &Widget{})
		assert.EqualError(t, fmt.Errorf("missing ID"), err.Error())
	})

	t.Run("Returns an error when the insert fails", func(t *testing.T) {
		// We never call setupTestWidgetsTable in this test
		Tx(t)
		dbWidgetMissingTable := newDBWidget()
		err := dbWidgetMissingTable.Create(ctx, &w1)
		assert.True(t, strings.Contains(err.Error(), "sqx_widgets_test' doesn't exist"))
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	status := "excellent"

	w1 := newWidget("great")

	t.Run("Can update a row as expected", func(t *testing.T) {
		setupTestWidgetsTable(t)
		dbWidget := newDBWidget()
		require.NoError(t, dbWidget.Create(ctx, &w1))
		assert.NoError(t, dbWidget.Update(ctx, w1.ID, &widgetUpdateFilter{
			Status: &status,
		}))
		w1db, err := dbWidget.GetByID(ctx, w1.ID)
		assert.NoError(t, err)
		expected := &Widget{
			ID:      w1.ID,
			Status:  "excellent",
			Enabled: w1.Enabled,
		}
		assert.Equal(t, expected, w1db)
	})

	t.Run("Can update a nullable row as expected", func(t *testing.T) {
		setupTestWidgetsTable(t)
		dbWidget := newDBWidget()
		ownerID := "owner-id"
		require.NoError(t, dbWidget.Create(ctx, &w1))

		// Owner IDs are null by default. Set the owner ID to a non-null value
		assert.NoError(t, dbWidget.Update(ctx, w1.ID, &widgetUpdateFilter{
			OwnerID: sqx.NewNullable[string](ownerID),
		}))
		w1db, err := dbWidget.GetByID(ctx, w1.ID)
		assert.NoError(t, err)
		assert.Equal(t, &ownerID, w1db.OwnerID)

		// Now set it back to null
		assert.NoError(t, dbWidget.Update(ctx, w1.ID, &widgetUpdateFilter{
			OwnerID: sqx.NewNull[string](),
		}))
		w2db, err := dbWidget.GetByID(ctx, w1.ID)
		assert.NoError(t, err)
		assert.Nil(t, w2db.OwnerID)
	})

	t.Run("Does not return an error on no updates", func(t *testing.T) {
		setupTestWidgetsTable(t)
		dbWidget := newDBWidget()
		// Empty update should not work
		err := dbWidget.Update(ctx, w1.ID, &widgetUpdateFilter{})
		assert.NoError(t, err)
	})

	t.Run("Returns an error when SetMap fails", func(t *testing.T) {
		setupTestWidgetsTable(t)
		dbWidget := newDBWidget()
		// Empty update should not work
		err := dbWidget.Update(ctx, w1.ID, &widgetUpdateFilter{Status: sqx.Ptr("Greasy")})
		assert.EqualError(t, err, "widgets cannot be greasy")
	})

	t.Run("Returns an error when the update fails", func(t *testing.T) {
		// We never call setupTestWidgetsTable in this test
		Tx(t)
		enabled := false
		dbWidgetMissingTable := newDBWidget()
		err := dbWidgetMissingTable.Update(ctx, w1.ID, &widgetUpdateFilter{
			Enabled: &enabled,
		})

		// Full error message: "Table 'testSQX.sqx_widgets_test' doesn't exist",
		// The database name may be different in different environments - only check the table name
		assert.True(t, strings.Contains(err.Error(), "sqx_widgets_test' doesn't exist"))
	})
}
