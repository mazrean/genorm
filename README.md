# GenORM

example: https://github.com/mazrean/genorm-workspace

## CLI(Code Generator)
### Install
```bash
$ go install github.com/mazrean/genorm/cmd/genorm@latest
```

### Code Generate
```
$ genorm -help
Usage of ./genorm:
  -destination string
    	The destination file to write.
  -join-num int
    	The number of joins to generate. (default 5)
  -module string
    	The root module name to use.
  -package string
    	The root package name to use.
  -source string
    	The source file to parse.
  -version
    	If true, output version information.
```

### Config

example: https://github.com/mazrean/genorm-workspace/blob/main/workspace/genorm_conf.go

#### Table
```go
type User struct {
	ID       uuid.UUID `genorm:"id"`
	Name     string    `genorm:"name"`
	Password string    `genorm:"password"`
}

func (u *User) TableName() string {
	return "users"
}
```

#### Table with Reference
```go
type Message struct {
	ID        uuid.UUID `genorm:"id"`
	UserID    uuid.UUID `genorm:"user_id"`
	Content   string    `genorm:"content"`
	CreatedAt time.Time `genorm:"created_at"`
	User      genorm.Ref[User]
}

func (m *Message) TableName() string {
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

### Transaction
Ref: https://pkg.go.dev/database/sql#example-DB.BeginTx

### Insert
```go
affectedRows, err := genorm.Insert(orm.User()).
  Values(&orm.UserTable{
    ID:       uuid.New(),
    Name:     genorm.Wrap("user"),
    Password: genorm.Wrap("password"),
  }).
  Do(db)
```

### Select

```go
userValues, err := genorm.Select(orm.User()).
  Fields(user.Name, user.Password).
  Where(genorm.EqLit(user.IDExpr, userID)).
  Find(db)
```

```go
userValue, err := genorm.Select(orm.User()).
  Fields(user.Name, user.Password).
  Where(genorm.EqLit(user.IDExpr, userID)).
  Take(db)
```

### Update
```go
affectedRows, err = genorm.Update(orm.Message()).
  Set(
    genorm.AssignLit(message.Content, genorm.Wrap("hello world")),
    genorm.AssignLit(message.CreatedAt, genorm.Wrap(time.Now())),
  ).
  Where(genorm.EqLit(message.IDExpr, messageID1)).
  Do(db)
```


### Delete
```go
affectedRows, err = genorm.Delete(orm.Message()).
  Where(genorm.EqLit(message.UserIDExpr, userID)).
  Do(db)
```

### Join
#### Select
```go
userIDColumn := orm.MessageUserParseExpr(user.ID)
messageUserIDColumn := orm.MessageUserParseExpr(message.UserID)
messageUserTable := orm.Message().
  User().Join(genorm.Eq(userIDColumn, messageUserIDColumn)
messageUserValues, err := genorm.Select(messageUserTable).
  Where(genorm.EqLit(userIDColumn, userID).
  Find(db)
```

#### Update
```go
userIDColumn := orm.MessageUserParseExpr(user.ID)
messageUserIDColumn := orm.MessageUserParseExpr(message.UserID)
messageContent := orm.MessageUserParseExpr(message.Content)
messageUserValues, err := genorm.Update(orm.Message().
    User().Join(genorm.Eq(userIDColumn, messageUserIDColumn))).
  Set(genorm.AssignLit(messageContent, genorm.Wrap("hello world"))).
  Where(genorm.EqLit(userIDColumn, userID)).
  Do(db)
```

### Context
```go
ctx := context.Background()
affectedRows, err := genorm.Insert(orm.User()).
  Values(&orm.UserTable{
    ID:       uuid.New(),
    Name:     genorm.Wrap("user"),
    Password: genorm.Wrap("password"),
  }).
  DoCtx(ctx, db)
```

