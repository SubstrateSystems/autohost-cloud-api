build:
	go build -o bin/autohost-cloud-lite ./cmd/api

test:
	go test ./...

docker:
	docker build -t autohost-cloud-lite .

dev:
	docker compose up -d
	air -c .air.toml


migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down 1

migrate-version:
	go run ./cmd/migrate version
