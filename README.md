# uses:

go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest

## usefull commands:

`sqlc generate`
`goose postgres postgres://postgres:postgres@localhost:5432/gator up`
`goose postgres postgres://postgres:postgres@localhost:5432/gator down`
