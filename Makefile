run:
	sqlc generate && templ generate --watch --proxy="http://localhost:3000" --cmd="go run main.go"