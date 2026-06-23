include .env
export

migration:
	goose -dir migrations create $(name) sql

migrate-up:
	goose -dir migrations postgres "$(PG_DSN)" up

migrate-down:
	goose -dir migrations postgres "$(PG_DSN)" down

migrate-status:
	goose -dir migrations postgres "$(PG_DSN)" status
