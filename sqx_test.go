package sqx_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dbConfig struct {
	ServerName string `env:"TEST_DB_SERVER_NAME"`
	User       string `env:"TEST_DB_USER"`
	Password   string `env:"TEST_DB_PASSWORD"`
	Name       string `env:"TEST_DB_NAME"`
}

func loadDBConfig(t *testing.T) (*dbConfig, error) {
	conf := new(dbConfig)
	var missing []string

	// conf must be a pointer type for this to be able to set values.
	vv := reflect.ValueOf(conf).Elem()
	tt := vv.Type()
	for i := 0; i < tt.NumField(); i++ {
		f := tt.Field(i)

		name := f.Tag.Get("env")
		if name == "" {
			continue
		}

		dest := vv.FieldByIndex(f.Index)
		val := os.Getenv(name)
		if val == "" {
			missing = append(missing, name)
		}
		dest.SetString(val)
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}
	return conf, nil
}

func CreateDatabase(c *dbConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		c.User,
		c.Password,
		c.ServerName,
		c.Name,
	)
	return sql.Open("mysql", connectionString)
}

// DB opens a new database connection for the duration of the test.
//
// If you don't need to persist data (e.g., for chained-request testing), use Tx instead.
func DB(t *testing.T) *sql.DB {
	conf, err := loadDBConfig(t)
	if err != nil {
		t.Fatalf("Open database connection: %s", err.Error())
	}

	db, err := CreateDatabase(conf)
	if err != nil {
		t.Fatalf("Create database: %s", err.Error())
	}

	t.Cleanup(func() {
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
		t.Fatalf("DB connection failed: %s", err.Error())
	}
	t.Cleanup(func() {
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
	_, err = db.Exec(`CREATE TABLE sqx_widgets_test (
    	widget_id      VARCHAR(128),
  		Status  VARCHAR(128),
  		enabled BOOLEAN
    )`)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := db.Exec(`DROP TABLE IF EXISTS sqx_widgets_test;`)
		require.NoError(t, err)
	})
	return db
}

var w1 = Widget{
	ID:      "widget_1",
	Status:  "great",
	Enabled: true,
}

var w2 = Widget{
	ID:      "widget_1",
	Status:  "fine",
	Enabled: true,
}

func TestRead(t *testing.T) {
	dbWidget := newDBWidget(setupTestWidgetsTable(t))
	ctx := context.Background()

	assert.NoError(t, dbWidget.Create(ctx, &w1))
	assert.NoError(t, dbWidget.Create(ctx, &w2))

	t.Run("Can read a single widget", func(t *testing.T) {
		w1db, err := dbWidget.GetByID(ctx, w1.ID)
		assert.NoError(t, err)
		assert.Equal(t, &w1, w1db)
	})

	t.Run("Can read multiple widgets", func(t *testing.T) {
		widgets, err := dbWidget.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []Widget{w1, w2}, widgets)
	})

	t.Run("Can read multiple widgets using a filter", func(t *testing.T) {
		widgets, err := dbWidget.Get(ctx, &widgetGetFilter{
			widgetID: &[]string{w1.ID, w2.ID},
		})
		assert.NoError(t, err)
		assert.Equal(t, []Widget{w1, w2}, widgets)
	})

	t.Run("Bubbles up errors from the DB exec", func(t *testing.T) {
		w1db, err := dbWidget.GetByID(ctx, "bad-id")
		assert.EqualError(t, sql.ErrNoRows, err.Error())
		assert.Nil(t, w1db)
	})
}

func TestInsert(t *testing.T) {
	ctx := context.Background()

	// The happy path cases are already tested in TestRead
	// here are only failure cases!

	t.Run("Returns an error when SetMap fails", func(t *testing.T) {
		dbWidget := newDBWidget(setupTestWidgetsTable(t))
		// Creating an empty widget should not work
		err := dbWidget.Create(ctx, &Widget{})
		assert.EqualError(t, fmt.Errorf("missing ID"), err.Error())
	})

	t.Run("Returns an error when the insert fails", func(t *testing.T) {
		// We never call setupTestWidgetsTable in this test
		db := Tx(t)
		dbWidgetMissingTable := newDBWidget(db)
		err := dbWidgetMissingTable.Create(ctx, &w1)
		assert.True(t, strings.Contains(err.Error(), "sqx_widgets_test' doesn't exist"))
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	status := "excellent"

	t.Run("Can update a row as expected", func(t *testing.T) {
		dbWidget := newDBWidget(setupTestWidgetsTable(t))
		assert.NoError(t, dbWidget.Create(ctx, &w1))
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

	t.Run("Does not return an error on no updates", func(t *testing.T) {
		dbWidget := newDBWidget(setupTestWidgetsTable(t))
		// Empty update should not work
		err := dbWidget.Update(ctx, w1.ID, &widgetUpdateFilter{})
		assert.NoError(t, err)
	})

	t.Run("Returns an error when SetMap fails", func(t *testing.T) {
		dbWidget := newDBWidget(setupTestWidgetsTable(t))
		// Empty update should not work
		err := dbWidget.Update(ctx, w1.ID, &widgetUpdateFilter{Status: ptr("Greasy")})
		assert.EqualError(t, err, "widgets cannot be greasy")
	})

	t.Run("Returns an error when the update fails", func(t *testing.T) {
		// We never call setupTestWidgetsTable in this test
		db := Tx(t)
		enabled := false
		dbWidgetMissingTable := newDBWidget(db)
		err := dbWidgetMissingTable.Update(ctx, w1.ID, &widgetUpdateFilter{
			Enabled: &enabled,
		})

		// Full error message: "Table 'testSQX.sqx_widgets_test' doesn't exist",
		// The database name may be different in different environments - only check the table name
		assert.True(t, strings.Contains(err.Error(), "sqx_widgets_test' doesn't exist"))
	})
}