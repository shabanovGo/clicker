.PHONY: up down migrate postgres recreate-db build logs test test-verbose test-coverage proto

DC=docker compose
DB_USER=clicks_user
DB_PASS=clicks_password
DB_NAME=clicks_db
DB_HOST=localhost
DB_PORT=5432
MIGRATIONS_DIR=migrations
PROJECT_NAME=clicks-counter
NETWORK=$(PROJECT_NAME)_clicks-network
PROTO_DIR=api/proto
GEN_DIR=pkg/clicker

up:
	$(DC) up

down:
	$(DC) down -v

build:
	$(DC) build

rebuild:
	$(DC) up --build

logs:
	$(DC) logs -f

postgres:
	$(DC) up -d postgres
	@until docker exec $$(docker ps -q -f name=postgres) pg_isready -U $(DB_USER) -d $(DB_NAME); do \
		echo "Waiting for postgres..."; \
		sleep 1; \
	done

recreate-db: postgres
	docker exec $$(docker ps -q -f name=postgres) dropdb -U $(DB_USER) --if-exists $(DB_NAME)
	docker exec $$(docker ps -q -f name=postgres) createdb -U $(DB_USER) $(DB_NAME)

migrate: recreate-db
	@for file in $(MIGRATIONS_DIR)/*.up.sql; do \
		echo "Applying $$file..."; \
		docker exec -i $$(docker ps -q -f name=postgres) psql -U $(DB_USER) -d $(DB_NAME) < $$file; \
	done

seed:
	@echo "Seeding database with sample banners..."
	docker exec -i $$(docker ps -q -f name=postgres) psql -U $(DB_USER) -d $(DB_NAME) < scripts/seed.sql
	@echo "Database seeded successfully!"

restart-app:
	$(DC) restart app

app-logs:
	$(DC) logs -f app

db-logs:
	$(DC) logs -f postgres

start: migrate seed build up

reset-db: recreate-db migrate seed

redis-cli:
	docker exec -it $$(docker ps -q -f name=redis) redis-cli

redis-logs:
	$(DC) logs -f redis

proto:
	@echo "Generating proto files..."
	@mkdir -p $(GEN_DIR)
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

.DEFAULT_GOAL := start
