# potom backend

## server

### install

```bash
go mod tidy
```

### run

```bash
go run .
./backend
```

### build

```bash
go build .

```

## db

### install

```bash
# install
brew install postgresql@15

# check version
psql --version

# start server in background
brew services start postgresql@15
```

### create table

```bash
psql postgres

# create db
CREATE DATABASE potom;

# connect to db
\c potom
```

### connection

connection string `protocol://username:password@host:port/database`

```bash
# don't need password on mac
psql postgres://moe:@localhost:5432/potom
psql potom
```

### migrations

install goose

> Goose is a database migration tool written in Go. It runs migrations from a set of SQL files, we wanna stay close to the raw SQL.

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

migrate

```bash
goose -dir sql/schema postgres postgres://moe:@localhost:5432/potom up
```

check migration

```bash
psql potom
\dt
```

### generate sql

> SQLC is a Go program that generates Go code from SQL queries. It's not exactly an ORM, but rather a tool that makes working with raw SQL easy and type-safe.

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc version
```

## swagger

```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

```bash
swag init
```