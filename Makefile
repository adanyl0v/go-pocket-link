include .env

POSTGRES_MIGRATIONS=./database/migrations/postgres
POSTGRES_CONNECTION=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_SSL_MODE}

all: build up migrate_up

local: migrate_up
	go run ./cmd/app/main.go

migrate_up:
	docker-compose up -d postgres
	goose -dir $(POSTGRES_MIGRATIONS) postgres $(POSTGRES_CONNECTION) up
	goose -dir $(POSTGRES_MIGRATIONS) postgres $(POSTGRES_CONNECTION) status

migrate_down:
	goose -dir $(POSTGRES_MIGRATIONS) postgres $(POSTGRES_CONNECTION) down

up:
	docker-compose up -d

down:
	docker-compose down

stop:
	docker-compose stop

build:
	docker-compose build

.SILENT: all local migrate_up migrate_down up down stop build
