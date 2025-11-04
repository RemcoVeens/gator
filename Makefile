reset_db:
	cd ./sql/schema/
	goose postgres postgres://postgres:postgres@localhost:5432/gator down --dir sql/schema
	goose postgres postgres://postgres:postgres@localhost:5432/gator up --dir sql/schema
