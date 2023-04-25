run:
	go run cmd/main.go

up:
	docker compose up -d

down:
	docker compose down

migrate-init:
	go run cmd/migrations/main.go db init

migrate:
	go run cmd/migrations/main.go db migrate

rollback:
	go run cmd/migrations/main.go db rollback

create-sql:
	go run cmd/migrations/main.go db create_sql