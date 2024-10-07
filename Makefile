include .env

DB_MIGRATIONS_DIR=./storage/migrations

DB_POSTGRES_DSN="postgres://root:${DB_POSTGRES_PASS}@localhost:5432/go_pocket_link?sslmode=disable"

all:
	@echo "\n- [+] Applying migrations..."
	goose -dir $(DB_MIGRATIONS_DIR) postgres $(DB_POSTGRES_DSN) status
	goose -dir $(DB_MIGRATIONS_DIR) postgres $(DB_POSTGRES_DSN) redo
	goose -dir $(DB_MIGRATIONS_DIR) postgres $(DB_POSTGRES_DSN) up
	@echo "\n- [+] Running the application..."
	go run ./cmd/app/main.go

up:
	@echo "\n- [+] Running docker containers"
	docker-compose up -d

up_build:
	@echo "\n- [+] Running docker containers after build"
	docker-compose up --build -d

down:
	@echo "\n- [+] Stopping docker containers"
	docker-compose down

.SILENT: all up up_build down
