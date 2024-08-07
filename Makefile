build:
	CGO_ENABLED=0 go build -o build/binapp cmd/main.go && docker compose build
