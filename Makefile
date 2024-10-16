include .env

POSTGRES_MIGRATIONS_PATH=./storage/postgres/migrations
POSTGRES_DSN="postgres://${POSTGRES_USER}:${POSTGRES_PASS}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_NAME}?sslmode=disable"

APP_EXECUTABLE=./.bin/go-pocket-link

all: deploy up

local:
	@echo "\n- [+] Applying migrations..."
	goose -dir $(POSTGRES_MIGRATIONS_PATH) postgres $(POSTGRES_DSN) up
	@echo "\n- [+] Building the application..."
	go build -o $(APP_EXECUTABLE) ./cmd/app/main.go
	@echo "- [+] Running the application..."
	$(APP_EXECUTABLE)

deploy:
	@echo "\n- [+] Building Docker containers..."
	docker-compose build

up:
	@echo "\n- [+] Running Docker containers..."
	docker-compose up -d

stop:
	@echo "\n- [+] Stopping Docker containers..."
	docker-compose stop

down:
	@echo "\n- [+] Deleting Docker containers..."
	docker-compose down

.SILENT: local deploy up stop down
