# Squirrel Xtended 
Squirrel Xtended (`sqx`) is a convenient library for db interactions in go. It provides nice bindings around:
- [Masterminds/squirrel](Masterminds/squirrel) - for fluent SQL generation
- [blockloop/scan](https://github.com/blockloop/scan) - for data marshalling
- [Go 1.18 Generics](https://go.dev/doc/tutorial/generics)

`sqx` is not an ORM or a migration tool. `sqx` just wants to run some SQL! 

### Example

```golang
package main

import (
	"context"
	sqx "github.com/stytchauth/sqx"
)

type User struct {
	ID          string `db:"id"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
	Status      string `db:"status"`
}

type GetUserFilter struct {
	ID          *string `db:"id"`
	Email       *string `db:"email"`
	PhoneNumber *string `db:"phone_number"`
}

func GetUsers(ctx context.Context, filter GetUserFilter) ([]User, error) {
	return sqx.Read[User]().
		Select("*").
		From("users").
		Where(sqx.ToClause(filter)).
		All(ctx)
}
```

---

### Setup
Teach `sqx` where your DB handle and logger are.
```golang

import sqx "github.com/stytchauth/sqx"

func main() {
	db := getDatabase()
	log := getLogger()
    sqx.SetDatabase(db)
	sqx.SetLogger(log)
}
```

Have multiple database handles or loggers? You can override them later.
```golang
func GetUsersFromReadReplica(ctx context.Context, filter GetUserFilter) ([]User, error) {
	return sqx.Read[User]().
		WithDatabase(replicaDB)
		Select("*").
		From("users").
		Where(sqx.ToClause(filter)).
		All(ctx)
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