# GenORM

SQL Builder to prevent SQL mistakes using Go language generics

#### document
- [English](https://mazrean.github.io/genorm-docs/en/)
- [日本語](https://mazrean.github.io/genorm-docs/ja/)

### Feature

By mapping SQL expressions to appropriate golang types using generics

* Compilation errors occur when using values of different Go types in SQL, such as for `=` comparisons or UPDATE value updates.
* Compilation error when using column names from a table that cannot be used

Many SQL mistakes can be discovered at the time of compilation that are not prevented by traditional Go language ORMs or query builders, such as

It also supports many CRUD syntaxes in SQL.

## Example
#### Example 1

String column `` `users`. `name` `` can be compared to a `string` value, but comparing it to an `int` value will result in a compile error.

```go
// correct
userValues, err := genorm.
	Select(orm.User()).
	Where(genorm.EqLit(user.NameExpr, genorm.Wrap("name"))).
	GetAll(db)

// compile error
userValues, err := genorm.
	Select(orm.User()).
	Where(genorm.EqLit(user.NameExpr, genorm.Wrap(1))).
	GetAll(db)
```

#### Example 2

You can use an `id` column from the `users` table in a `SELECT` statement that retrieves data from the `users` table, but using an `id` column from the `messages` table will result in a compile error.

```go
// correct
userValues, err := genorm.
	Select(orm.User()).
	Where(genorm.EqLit(user.IDExpr, uuid.New())).
	GetAll(db)

// compile error
userValues, err := genorm.
	Select(orm.User()).
	Where(genorm.EqLit(message.IDExpr, uuid.New())).
	GetAll(db)
```

## Install

GenORM uses the CLI to generate code. The `genorm`package is used to invoke queries. For this reason, both the CLI and Package must be install.

#### CLI

```
go install github.com/mazrean/genorm/cmd/genorm@v1.0.0
```

#### Package

```
go get -u github.com/mazrean/genorm
```

### Configuration

#### Example

The `users` table can join the `messages` table.

```go
import "github.com/mazrean/genorm"

type User struct {
    // Column Information
    Message genorm.Ref[Message]
}

func (*User) TableName() string {
    return "users"
}

type Message struct {
    // Column Information
}

func (*Message) TableName() string {
    return "messages"
}
```

## Usage
### Connecting to a Database
```go
import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

db, err := sql.Open("mysql", "user:pass@tcp(host:port)/database?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4")
```

### Insert
```go
// INSERT INTO `users` (`id`, `name`, `created_at`) VALUES ({{uuid.New()}}, "name1", {{time.Now()}}), ({{uuid.New()}}, "name2", {{time.Now()}})
affectedRows, err := genorm.
    Insert(orm.User()).
    Values(&orm.UserTable{
        ID: uuid.New(),
        Name: genorm.Wrap("name1"),
        CreatedAt: genorm.Wrap(time.Now()),
    }, &orm.UserTable{
        ID: uuid.New(),
        Name: genorm.Wrap("name2"),
        CreatedAt: genorm.Wrap(time.Now()),
    }).
    Do(db)
```

### Select

```go
// SELECT `id`, `name`, `created_at` FROM `users`
// userValues: []orm.UserTable
userValues, err := genorm.
	Select(orm.User()).
	GetAll(db)

// SELECT `id`, `name`, `created_at` FROM `users` LIMIT 1
// userValue: orm.UserTable
userValue, err := genorm.
	Select(orm.User()).
	Get(db)

// SELECT `id` FROM `users`
// userIDs: []uuid.UUID
userIDs, err := genorm.
	Pluck(orm.User(), user.IDExpr).
	GetAll(db)

// SELECT COUNT(`id`) AS `result` FROM `users` LIMIT 1
// userNum: int64
userNum, err := genorm.
	Pluck(orm.User(), genorm.Count(user.IDExpr, false)).
	Get(db)
```

### Update
```go
// UPDATE `users` SET `name`="name"
affectedRows, err = genorm.
    Update(orm.User()).
    Set(
        genorm.AssignLit(user.Name, genorm.Wrap("name")),
    ).
    Do(db)
```


### Delete
```go
// DELETE FROM `users`
affectedRows, err = genorm.
    Delete(orm.User()).
    Do(db)
```

### Join
#### Select
```go
// SELECT `users`.`name`, `messages`.`content` FROM `users` INNER JOIN `messages` ON `users`.`id` = `messages`.`user_id`
// messageUserValues: []orm.MessageUserTable
userID := orm.MessageUserParseExpr(user.ID)
userName := orm.MessageUserParse(user.Name)
messageUserID := orm.MessageUserParseExpr(message.UserID)
messageContent := orm.MessageUserParse(message.Content)
messageUserValues, err := genorm.
	Select(orm.User().
		Message().Join(genorm.Eq(userID, messageUserID))).
	Fields(userName, messageContent).
	GetAll(db)
```

#### Update
```go
// UPDATE `users` INNER JOIN `messages` ON `users.id` = `messages`.`id` SET `content`="hello world"
userIDColumn := orm.MessageUserParseExpr(user.ID)
messageUserIDColumn := orm.MessageUserParseExpr(message.UserID)
messageContent := orm.MessageUserParse(message.Content)
affectedRows, err := genorm.
  Update(orm.User().
		Message().Join(genorm.Eq(userID, messageUserID))).
  Set(genorm.AssignLit(messageContent, genorm.Wrap("hello world"))).
  Do(db)
```

### Transaction
```go
tx, err := db.Begin()
if err != nil {
    log.Fatal(err)
}

_, err = genorm.
    Insert(orm.User()).
    Values(&orm.UserTable{
        ID: uuid.New(),
        Name: genorm.Wrap("name1"),
        CreatedAt: genorm.Wrap(time.Now()),
    }, &orm.UserTable{
        ID: uuid.New(),
        Name: genorm.Wrap("name2"),
        CreatedAt: genorm.Wrap(time.Now()),
    }).
    Do(db)
if err != nil {
    _ = tx.Rollback()
    log.Fatal(err)
}

err = tx.Commit()
if err != nil {
    log.Fatal(err)
}
```

### Context

```go
// SELECT `id`, `name`, `created_at` FROM `users`
// userValues: []orm.UserTable
userValues, err := genorm.
	Select(orm.User()).
	GetAllCtx(context.Background(), db)
```

```go
// INSERT INTO `users` (`id`, `name`, `created_at`) VALUES ({{uuid.New()}}, "name", {{time.Now()}})
affectedRows, err := genorm.
    Insert(orm.User()).
    Values(&orm.UserTable{
        ID: uuid.New(),
        Name: genorm.Wrap("name"),
        CreatedAt: genorm.Wrap(time.Now()),
    }).
    DoCtx(context.Background(), db)
```

