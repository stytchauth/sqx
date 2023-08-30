# Squirrel Xtended 
Squirrel Xtended (`sqx`) is a convenient library for db interactions in go. It provides nice bindings around:
- [Masterminds/squirrel](https://github.com/Masterminds/squirrel) - for fluent SQL generation
- [blockloop/scan](https://github.com/blockloop/scan) - for data marshalling
- [Go 1.18 Generics](https://go.dev/doc/tutorial/generics)

`sqx` is not an ORM or a migration tool. `sqx` just wants to run some SQL! 

### Quick Start
Teach `sqx` where your DB handle and logger are. `sqx` can then be used to create, update, and delete data.

See [Widget Test](./widget_test.go) for an example of a complete data layer built with `sqx`.

```golang
package main

import (
	"context"
	"github.com/stytchauth/sqx"
)

func init() {
	db := getDatabase()
	log := getLogger()
	sqx.SetDefaultQueryable(db)
	sqx.SetDefaultLogger(log)
}

type User struct {
	ID          string `db:"id"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
	Status      string `db:"status"`
}

func InsertUser(ctx context.Context, user *User) error {
	return sqx.Write(ctx).
		Insert("users").
		SetMap(sqx.ToSetMap(user)).
		Do()
}

type GetUserFilter struct {
	ID          *string `db:"id"`
	Email       *string `db:"email"`
	PhoneNumber *string `db:"phone_number"`
}

func GetUsers(ctx context.Context, filter GetUserFilter) ([]User, error) {
	return sqx.Read[User](ctx).
		Select("*").
		From("users").
		Where(sqx.ToClause(filter)).
		All()
}

func DeleteUser(ctx context.Context, userID string) error {
	return sqx.Write(ctx).
		Delete("users").
		Where(sqx.Eq{"ID": userID}).
		Do()
}

```
---
### Core Concepts

#### Query building

`sqx` is a superset of [Masterminds/squirrel](https://github.com/Masterminds/squirrel) - refer to their docs for information on what query methods are available.
We will try to add more examples over time - if there is an example you'd love to see, feel free to open an issue or a PR!

#### Reading data
Call `sqx.Read[T](ctx).Select(columNames...)` to start building a read transaction. When the read transaction is ran, 
`sqx` will provision an object of type `T` and scan the results into the object. Scanning is accomplished using `db` tags defined on `T`. 
All scanning is handled by [blockloop/scan](https://github.com/blockloop/scan)'s `RowsStrict` method.
Read transactions can be ran in several ways:

- `func (b SelectBuilder[T]) One() (*T, error)` - reads a single struct of type `T`. 
If no response is found, returns a `sql.ErrNoRows`.
If more than one row is returned from the underlying query, an error will be logged to the provided logger.
- `func (b SelectBuilder[T]) OneStrict() (*T, error)` - like `One()` but returns an error if more than one row is returned
- `func (b SelectBuilder[T]) OneScalar() (T, error)` - like `One()` but can be used to read simple values like `int32` or `string`
- `func (b SelectBuilder[T]) First() (*T, error)` - line `One()` but does not care if the underlying query has more than
  one result and will just take the first row.
**NOTE**: if you don't supply an OrderBy clause, the first result is not guaranteed to be the same each time you run the
query.
- `func (b SelectBuilder[T]) FirstScalar() (T, error)` - line `First()` but can be used to read simple values like
  `int32` or `string`
- `func (b SelectBuilder[T]) All() ([]T, error)` - returns a slice of structs of type `T`

You'll often want to filter the data that you read - for example, finding all `Users` with a certain status, or finding a `User` with a specific ID.
`sqx.ToClause` is helpful for converting flexible structs into `Where`-compatible filters. `nil`-valued fields are ignored,
and only present fields are preserved.

For example, the following struct definition can be used to find users with a specific ID, a specific Email, a specific PhoneNumber, or any combination thereof.
```golang
type GetUserFilter struct {
	ID          *string `db:"id"`
	Email       *string `db:"email"`
	PhoneNumber *string `db:"phone_number"`
}
```

- `sqx.ToClause(GetUserFilter{ID: sqx.Ptr("123")})` -> `sqx.Eq{"id": "123"}`
- `sqx.ToClause(GetUserFilter{Email: sqx.Ptr("joe@example.com")})` -> `sqx.Eq{"email": "joe@example.com"}`
- `sqx.ToClause(GetUserFilter{ID: sqx.Ptr("123"), Email: sqx.Ptr("joe@example.com")})` -> `sqx.Eq{"id": "123", "email": "joe@example.com"}`

```golang
func GetUsers(ctx context.Context, filter GetUserFilter) ([]User, error) {
	return sqx.Read[User](ctx).
		Select("*").
		From("users").
		Where(sqx.ToClause(filter)).
		All()
}
```

If you are joining tables together and aliasing them along the way, `sqx.ToClauseAlias` can help with that.
```golang
func GetUsersAndProfileData(ctx context.Context, filter GetUserFilter) ([]User, error) {
	return sqx.Read[UserWithPets](ctx).
		Select("*").
		From("users u").
		Join("pets p ON users.id = pets.user_id")
		Where(sqx.ToClauseAlias("u", filter)).
		All()
}
```

You can also define the alias directly in the struct tag
```golang
type GetUserWithPetFilter struct {
	UserID *string `db:"u.id"`
	PetID  *string `db:"p.id"`
}
```

#### Writing data
Call `sqx.Write(ctx)` to start building a write transaction. Write transactions can be used for `Create`, `Update`, or `Delete` operations.
All write transactions are ran by calling `.Do()` after being built.

Create and Update transactions require fields to be set. Fields may be set one at a time via calls to `.Set(fieldName string, fieldValue any)` but the preferred way
is via `.SetMap(map[string]any)`. The method `sqx.ToSetMap` is useful for converting flexible structs into maps. 
As with `ToClause`, `nil`-valued fields are ignored, and only present fields are passed through.

For example, the following structs define a user that can be created once, then updated any number of times.
The `UserUpdate` struct can be used to update a user's email, phone number, status, or multiple at once.
```golang
type User struct {
	ID          string `db:"id"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
	Status      string `db:"status"`
}
type UserUpdate struct {
	Email       *string `db:"email"`
	PhoneNumber *string `db:"phone_number"`
	Status      *string `db:"status"`
}
```

- `sqx.ToSetMap(User{ID:"123", Email:"joe@example.com"})` -> `map[string]any{"id":"123", "email":"joe@example.com", "phone_number": "", "status":""}`
- `sqx.ToSetMap(UserUpdate{ID:sqx.Ptr("123"), Email:sqx.Email("joe@example.com")})` -> `map[string]any{"id":"123", "email":"joe@example.com"}`

```golang
func InsertUser(ctx context.Context, user *User) error {
	return sqx.Write(ctx).
		Insert("users").
		SetMap(sqx.ToSetMap(user)).
		Do()
}

func UpdateUser(ctx context.Context, userID string, update *UserUpdate) error {
	return sqx.Write(ctx).
		Update("users").
		Where(sqx.Eq{"id": userID}).
		SetMap(sqx.ToSetMap(update)).
		Do()
}
```

---

### Examples

#### Reading a single struct row
```golang
func GetUser(ctx context.Context, userID string) (*User, error) {
	return sqx.Read[User](ctx).
		Select("*").
		From("users").
		Where(sqx.Eq{"ID": userID}).
		One()
}
```

#### Reading a simple value (string, int, bool, etc)
```golang
func CountUsers(ctx context.Context, userID string) (int32, error) {
	return sqx.Read[int32](ctx).
		Select("COUNT(*)").
		From("users").
		OneScalar()
}
```

#### Reading a slice of structs
```golang
func GetAllUsers(ctx context.Context) ([]User, error) {
	return sqx.Read[User](ctx).
		Select("*").
		From("users").
		All()
}
```

#### Debugging generated SQL
Call `.Debug()` at any time to print out the internal state of the query builder
```golang
sqx.Read[UserWithPets](ctx).
	Select("*").
	From("users u").
	Debug().
	Join("pets p ON users.id = pets.user_id").
	Where(sqx.ToClauseAlias("u", filter)).
	Debug().
	All()
// outputs
// map[args:[] error:<nil> sql:SELECT * FROM users u]
// map[args:[poodle] error:<nil> sql:SELECT * FROM users u JOIN pets p ON users.id = pets.user_id WHERE u.breed = ?]
```

#### Setting a field to `null` using an Update
Use the `sqx.Nullable[T]` type and its helper methods - `sqx.NewNullable` and `sqx.NewNull`.

Given the update request:
```golang
type PetUpdate {
	UserID sqx.Nullable[string] `db:"user_id"`
}
func UpdatePets(ctx context.Context, petID string, petUpdate *PetUpdate) error {
	return sqx.Write(ctx).
		Update("pets").
		Where(sqx.Eq{"id": petID}).
		SetMap(sqx.ToClause(petUpdate)).
		Do()
}
```
This update will set the `user_id` field to the provided value
```golang
UpdatePets(ctx, &PetUpdate{
	UserID: sqx.NewNullable("some-user-id")
})
```
and this update will set the `user_id` field to `NULL`/`nil`
```golang
UpdatePets(ctx, &PetUpdate{
	UserID: sqx.NewNull[string]()
})
```

#### Validating data before inserting
`InsertBuilder.SetMap()` can take in an optional error. If an error occurs, the insert operation will short-circuit.

```golang
type Pet struct {
	Name string `db:"name"`
}
func (p *Pet) ToSetMap() (map[string]any, error) {
	if p.name == "" {
		return nil, fmt.Errorf("pet was missing name")		
	}
	return sqx.ToSetMap(p), nil
}

func CreatePet(ctx context.Context, pet *Pet) error {
	return sqx.Write(ctx).
		Insert("pets").
		SetMap(pet.ToSetMap()).
		Do()
}
```

#### Managing Transactions
`sqx` does not manage transactions itself. Create transactions within your application when needed, and then pass to
`WithQueryable` to let the request builder know to use that transaction object. Both `sql.DB` and `sql.Tx` satisfy the `sqx.Queryable` interface.

```golang
func MyOperationThatNeedsATransaction(ctx context.Context) error {
	// Get a Tx for making transaction requests.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()
	
	err = OperationThatNeedsAQueryable(ctx, tx);
	if err != nil {
		return err
	}

	err = OperationThatNeedsAQueryable(ctx, tx);
	if err != nil {
		return err
	}
  
	return tx.Commit()
}

func OperationThatNeedsAQueryable(ctx context.Context, tx sqx.Queryable) error {
	return sqx.Write(ctx).
		WithQueryable(tx).
		Update("table").
		Set("key", "value").
		Do()
}

```

#### Customizing Handles & Loggers

Have multiple database handles or a per-request logger? You can override them using `WithQueryable` or `WithLogger`.
```golang
func GetUsersFromReadReplica(ctx context.Context, filter GetUserFilter) ([]User, error) {
	return sqx.Read[User](ctx).
		WithQueryable(replicaDB).
		WithLogger(logging.FromCtx(ctx))
		Select("*").
		From("users").
		Where(sqx.ToClause(filter)).
		All()
}
```

If you always want to pass in a custom handle or logger, consider aliasing the `Read` and `Write` entrypoints within your project.

```golang
func Read[T any](ctx context.Context, db sqx.Queryable) interface {
	Select(columns ...string) sqx.SelectBuilder[T]
} {
	return sqx.Read[T](ctx).WithQueryable(db).WithLogger(logging.FromContext(ctx))
}

func Write(ctx context.Context, db sqx.Queryable) interface {
	Insert(tblName string) sqx.InsertBuilder
	Update(tblName string) sqx.UpdateBuilder
	Delete(tblName string) sqx.DeleteBuilder
} {
	return sqx.Write(ctx).WithQueryable(db).WithLogger(logging.FromContext(ctx))
}
```
----
### Contributing
`sqx` uses `mysql@8.1.0` in a docker file for development and testing. It is hardcoded to run on port `4306`

Start it with
```bash
make services
```

and kill it with
```bash
make services-stop
```

Run all tests with
```bash
make tests
```
