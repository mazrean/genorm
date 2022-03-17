# GenORM

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![](https://pkg.go.dev/badge/github.com/mazrean/genorm)](https://pkg.go.dev/github.com/mazrean/genorm)
[![](https://github.com/mazrean/genorm/workflows/CI/badge.svg)](https://github.com/mazrean/genorm/actions)

SQL Builder to prevent SQL mistakes using the Golang generics

#### document
- [English](https://mazrean.github.io/genorm-docs/en/)
- [日本語](https://mazrean.github.io/genorm-docs/ja/)

## Feature

By mapping SQL expressions to appropriate golang types using generics, you can discover many SQL mistakes at the time of compilation that are not prevented by traditional Golang ORMs or query builders.

For example:

* Compilation error occurs when using values of different Go types in SQL for = etc. comparisons or updating values in UPDATE statements
* Compile error occurs when using column names of unavailable tables

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

## Differences from existing tools
Explain how this differs from existing GORM and ent.

### GORM
GORM uses `interface{}` to flexibly accept arguments, allowing for intuitive query construction.
On the other hand, type constraints are very weak.
For this reason, there are few bugs that can be played at compile time, and bugs can easily be introduced.
Another problem is that it is difficult to understand the SQL that is actually executed, and it is easy for the SQL to behave differently from what is expected.

In contrast, GenORM uses type constraints to find as many bugs as possible at compile time.
We are also conscious of making it easy to imagine the SQL to be executed from the Go code, and we hope that the example code so far has been easy to imagine the SQL to be executed from the code.

Thus, GenORM and GORM differ greatly in both philosophy and function.

### ent
We believe that the most significant difference in functionality between Go and ent is the Golang types that can be used.
In ent, only a finite number of types can be used, including types corresponding to primitive types such as int and bool, plus time.Time and UUID.
In contrast, GenORM allows any type[^2] that satisfies the conditions of the `genorm.ExprType`interface, and then sets restrictions such as only types of the same type can be compared.
This eliminates the need for unnecessary type conversions and allows for stronger constraints to be set by using Defined Type[^3].

For example, by defining `UserID` and `MessageID` as unique types as shown below, constraints such as "user IDs and message IDs cannot be compared" can be specified.
```go
type MessageID uuid.UUID

func (mid *MessageID) Scan(src any) error {
	return (*uuid.UUID)(mid).Scan(src)
}

func (mid MessageID) Value() (driver.Value, error) {
	return uuid.UUID(mid).Value()
}

type UserID uuid.UUID

func (uid *UserID) Scan(src any) error {
	return (*uuid.UUID)(uid).Scan(src)
}

func (uid UserID) Value() (driver.Value, error) {
	return uuid.UUID(uid).Value()
}
```

This prevents mistakes, as code such as the following will result in a compilation error.
```go
// SELECT * FROM `messages` WHERE `messages`.`id`=`messages`.`user_id`
messageValues, err := genorm.
    Select(orm.Message()).
    Where(genorm.Eq(message.IDExpr, message.UserIDExpr)).
    GetAll(db)
```

It is also important to note that GenORM allows queries to be constructed with a method chain similar to SQL.
Because ent is an "entity framework," database operations are abstracted as entity operations.
This has its advantages, but as with GORM, it also has the aspect of making the SQL to be executed difficult to understand.
This increases the probability of unintentionally writing a process that has performance problems.
In this respect, GenORM makes it easy to understand the SQL to be executed, so such problems are less likely to occur.

[^2]: See https://mazrean.github.io/genorm-docs/en/usage/value-type.html for more information
[^3]: See https://mazrean.github.io/genorm-docs/en/advanced-usage/defined-type.html for more information

## Mechanism
This section explains how the code example works.
The definition of `genorm.EqLit` is as follows.
```go
func EqLit[T Table, S ExprType](
	expr TypedTableExpr[T, S],
	literal S,
) TypedTableExpr[T, WrappedPrimitive[bool]] {
	// 省略
}
```
Notice the `TypedTableExpr[T, S]` part.
As you can see from the `[T Table, S ExprType]` part, `T` is the type that represents the table used by the SQL expression, and `S` is the Go language type corresponding to the SQL expression.
Thus, GenORM limits the columns that can be used and the values that can be compared by having the table used by the SQL expression and the corresponding type of the Go language as type parameters.

Thus, because SQL expressions are typed, SQL can be restricted to operators such as `>`, `<`, and `AND`, functions such as `COUNT()`, and in the future, database-specific functions as well.

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

