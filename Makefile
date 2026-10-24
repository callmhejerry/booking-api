include .env
export

.PHONY: up down

up:
	docker compose up -d

down:
	docker compose down

db-shell:
	docker compose exec postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}


DB_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}

migrate-create:
	@read -p "Enter Migration name: " name; \
	migrate create -ext sql -dir sql/migrations -seq $$name

migrate-up:
	migrate -database "$(DB_URL)" -path sql/migrations up

migrate-down:
	migrate -database "$(DB_URL)" -path sql/migrations down

migrate-force:
	@read -p "Enter db version: " version; \
	migrate -database $(DB_URL) -path sql/migrations force $$version

migrate-version:
	migrate -database "$(DB_URL)" -path sql/migrations version


run:
	go run -race ./cmd/api/main.go

test:
	go test -v -count=1 ./tests/e2e/...