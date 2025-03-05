dev-start: 
	npx @tailwindcss/cli -i ./static/input.css -o ./static/output.css && templ generate && go run ./cmd/server/main.go
