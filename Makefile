.PHONY:
.SILENT:
.DEFAULT_GOAL := run

compose-up:
	docker-compose up --build && docker-compose logs -f

export PG_URL=postgres://postgres:postgres@localhost:5433/banner?sslmode=disable
run:
	go mod tidy && go mod download && go run -tags migrate ./cmd/app

export TEST_DB_URL=postgres://postgres:postgres@localhost:5432/banner?sslmode=disable
test:
	docker run --rm -d -p 5432:5432 --name test_db -e "POSTGRES_DB=banner" -e "POSTGRES_USER=postgres" -e "POSTGRES_PASSWORD=postgres" -e "POSTGRES_HOST=postgres-db" -e "POSTGRES_PORT=5432" -e "POSTGRES_SSLMODE=disable" postgres:latest
	go test -v ./tests/
	docker stop test_db

lint:
	golangci-lint run