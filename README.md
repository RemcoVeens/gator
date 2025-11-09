# gator

## dependecies:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## how to install:

```bash
go install www.github.com/RemcoVeens/gator.git@latest
```

## how to use:

1. create a `.gatorconfig.json` in your home dir
1. add the following data to it, and customize it to your system

```json
{
  "db_url": "postgres://<username>:<password>@localhost:5432/gator?sslmode=disable"
}
```

3. set up database by running `make db_up`
1. create a user by running `gator register <your_name>`
1. login useing `gator login <your_name>`
1. add a feed using `gator addfeed "feedname" "www.feed.com"`
1. farm data using `gator agg`
1. read articals using `gator browse`

## dev notes:

### usefull commands:

`sqlc generate` adds new quieries to database moduale
