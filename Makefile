include .env

POSTGRES_MIGRATIONS_PATH=./storage/postgres/migrations
POSTGRES_DSN="postgres://${POSTGRES_USER}:${POSTGRES_PASS}@localhost:5432/go_pocket_link?sslmode=disable"

APP_EXECUTABLE=./build/go-pocket-link

all:
	@echo "\n- [+] Applying migrations..."
	goose -dir $(POSTGRES_MIGRATIONS_PATH) postgres $(POSTGRES_DSN) up
	@echo "\n- [+] Building the application..."
	go build -o $(APP_EXECUTABLE) ./cmd/app/main.go
	@echo "- [+] Running the application..."
	$(APP_EXECUTABLE)

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
