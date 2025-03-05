dev-start: setup
	go run ./cmd/server/main.go

setup:
	npx @tailwindcss/cli -i ./static/input.css -o ./static/output.css --minify && node build.js && templ generate
