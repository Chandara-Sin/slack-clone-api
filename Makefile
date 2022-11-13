run:
	go run cmd/main.go

up:
	docker compose up -d

down:
	docker compose down

migrate:
	go run cmd/migrations/main.go db migrate

migrate-init:
	go run cmd/migrations/main.go db init

create-sql:
	go run cmd/migrations/main.go db create_sql