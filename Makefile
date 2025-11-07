db_down:
	goose postgres postgres://postgres:postgres@localhost:5432/gator down --dir sql/schema
db_up:
	goose postgres postgres://postgres:postgres@localhost:5432/gator up --dir sql/schema

reset_db:
	make db_down
	make db_up
