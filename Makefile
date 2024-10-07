include .env

DB_MIGRATIONS_DIR=./storage/migrations
DB_POSTGRES_DSN="postgres://${DB_POSTGRES_USER}:${DB_POSTGRES_PASS}@${DB_POSTGRES_HOST}:${DB_POSTGRES_PORT}/${DB_POSTGRES_NAME}?sslmode=${DB_POSTGRES_SSL_MODE}"

all:
	@echo "\n- [+] Applying migrations..."
	goose -dir $(DB_MIGRATIONS_DIR) postgres $(DB_POSTGRES_DSN) status
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
